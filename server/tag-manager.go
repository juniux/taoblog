package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/movsb/taoblog/modules/sql_helpers"
)

// TagNotFoundError is
type TagNotFoundError struct {
}

func (z *TagNotFoundError) Error() string {
	return "tag not found"
}

// Tag is a tag.
type Tag struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Alias int64  `json:"alias"`
}

// TagWithCount is a tag with associated post count.
type TagWithCount struct {
	Tag
	Count int64 `json:"count"`
}

// TagManager manages tags.
type TagManager struct {
}

// NewTagManager news
func NewTagManager() *TagManager {
	return &TagManager{}
}

func (tm *TagManager) addTag(tx Querier, name string, alias uint) int64 {
	sql := `INSERT INTO tags (name,alias) values (?,?)`
	ret, err := tx.Exec(sql, name, alias)
	if err != nil {
		panic(err)
	}

	id, err := ret.LastInsertId()
	return id
}

func (tm *TagManager) searchTag(tx Querier, tag string) (tags []Tag) {
	sql := `SELECT * FROM tags WHERE name LIKE ?`
	rows, err := tx.Query(sql, "%"+tag+"%")
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var tag Tag
		if err = rows.Scan(&tag.ID, &tag.Name, &tag.Alias); err != nil {
			panic(err)
		}
		tags = append(tags, tag)
	}

	return
}

func (tm *TagManager) getTagID(tx Querier, name string) int64 {
	sql := `SELECT id FROM tags WHERE name=? LIMIT 1`
	row := tx.QueryRow(sql, name)

	var id int64
	if err := row.Scan(&id); err != nil {
		log.Printf("标签名不存在：%s\n", name)
		id = 0
	}

	return id
}

func (tm *TagManager) hasTagName(tx Querier, name string) bool {
	return tm.getTagID(tx, name) > 0
}

// GetObjectTagNames gets all tag names of an object.
func (tm *TagManager) GetObjectTagNames(tx Querier, oid int64) ([]string, error) {
	query, args := sql_helpers.NewSelect().From("post_tags", "").From("tags", "").
		Select("tags.name").Where("post_tags.post_id=?", oid).Where("post_tags.tag_id=tags.id").SQL()
	rows, err := tx.Query(query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	names := make([]string, 0)

	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			return nil, err
		}
		names = append(names, name)
	}

	return names, rows.Err()
}

func (tm *TagManager) getTagIDs(tx Querier, pid int64, alias bool) (ids []int64) {
	sql := `SELECT tag_id FROM post_tags WHERE post_id=` + fmt.Sprint(pid)

	rows, err := tx.Query(sql)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		var id int64
		err = rows.Scan(&id)
		if err != nil {
			panic(err)
		}
		ids = append(ids, id)
	}

	if alias {
		ids = tm.getAliasTagsAll(tx, ids)
	}

	return
}

func (tm *TagManager) getAliasTagsAll(tx Querier, ids []int64) []int64 {
	sids := joinInts(ids, ",")

	sql1 := `SELECT alias FROM tags WHERE id in (?)`
	sql2 := `SELECT id FROM tags WHERE alias in (?)`

	rows, err := tx.Query(sql1, sids)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var alias int64
		if err = rows.Scan(&alias); err != nil {
			panic(err)
		}

		if alias > 0 {
			ids = append(ids, alias)
		}
	}

	rows.Close()

	rows, err = tx.Query(sql2, sids)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var id int64
		if err = rows.Scan(&id); err != nil {
			panic(err)
		}

		ids = append(ids, id)
	}

	rows.Close()

	return ids
}

func (tm *TagManager) addObjectTag(tx Querier, pid int64, tid int64) int64 {
	sql := `INSERT INTO post_tags (post_id,tag_id) VALUES (?,?)`
	ret, err := tx.Exec(sql, pid, tid)
	if err != nil {
		panic(err)
	}

	id, err := ret.LastInsertId()

	return id
}

func (tm *TagManager) removeObjectTag(tx Querier, pid, tid int64) {
	sql := `DELETE FROM post_tags WHERE post_id=? AND tag_id=? LIMIT 1`
	ret, err := tx.Exec(sql, pid, tid)
	if err != nil {
		panic(err)
	}
	_ = ret
}

// this is temp
func (tm *TagManager) hasObjectTag(tx Querier, pid, tid int64) bool {
	query := "SELECT id FROM post_tags WHERE post_id=? AND tag_id=? LIMIT 1"
	row := tx.QueryRow(query, pid, tid)
	id := 0
	row.Scan(&id)
	return id > 0
}

// UpdateObjectTags updates
func (tm *TagManager) UpdateObjectTags(tx Querier, pid int64, tags []string) error {
	newTags := tags
	oldTags, err := tm.GetObjectTagNames(tx, pid)
	if err != nil {
		return err
	}

	var (
		toBeDeled []string
		toBeAdded []string
	)

	for _, t := range oldTags {
		if !strInSlice(newTags, t) {
			toBeDeled = append(toBeDeled, t)
		}
	}

	for _, t := range newTags {
		t = strings.TrimSpace(t)
		if t != "" && !strInSlice(oldTags, t) {
			toBeAdded = append(toBeAdded, t)
		}
	}

	for _, t := range toBeDeled {
		tid := tm.getTagID(tx, t)
		tm.removeObjectTag(tx, pid, tid)
	}

	for _, t := range toBeAdded {
		var tid int64
		if !tm.hasTagName(tx, t) {
			tid = tm.addTag(tx, t, 0)
		} else {
			tag, _ := tm.GetRootTag(tx, tm.getTagID(tx, t))
			log.Println("root tag:", tag)
			tid = tag.ID
		}
		if !tm.hasObjectTag(tx, pid, tid) {
			tm.addObjectTag(tx, pid, tid)
		}
	}

	return nil
}

// ListTags lists all tags.
func (tm *TagManager) ListTags(tx Querier) ([]*Tag, error) {
	rows, err := tx.Query("SELECT * FROM tags")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	tags := make([]*Tag, 0)

	for rows.Next() {
		var tag Tag
		if err = rows.Scan(&tag.ID, &tag.Name, &tag.Alias); err != nil {
			return nil, err
		}
		tags = append(tags, &tag)
	}

	return tags, rows.Err()
}

// GetTagByID gets a tag by ID.
func (tm *TagManager) GetTagByID(tx Querier, id int64) (*Tag, error) {
	query := "SELECT id,name,alias FROM tags WHERE id=?"
	row := tx.QueryRow(query, id)
	var tag Tag
	if err := row.Scan(&tag.ID, &tag.Name, &tag.Alias); err != nil {
		if err == sql.ErrNoRows {
			return nil, &TagNotFoundError{}
		}
		return nil, err
	}
	return &tag, nil
}

// GetTagByName gets a tag by Name.
func (tm *TagManager) GetTagByName(tx Querier, tag string) (*Tag, error) {
	query := "SELECT id,name,alias FROM tags WHERE name=?"
	row := tx.QueryRow(query, tag)
	var tagObj Tag
	if err := row.Scan(&tagObj.ID, &tagObj.Name, &tagObj.Alias); err != nil {
		if err == sql.ErrNoRows {
			return nil, &TagNotFoundError{}
		}
		return nil, err
	}
	return &tagObj, nil
}

// GetRootTag gets the root tag of an alias-ed tag.
func (tm *TagManager) GetRootTag(tx Querier, id int64) (tag *Tag, err error) {
	for {
		tag, err = tm.GetTagByID(tx, id)
		if err != nil {
			return
		}
		if tag.Alias == 0 {
			return
		}
		id = tag.Alias
	}
}

// UpdateTag updates a tag.
func (tm *TagManager) UpdateTag(tx Querier, tag *Tag) error {
	if _, err := tm.GetTagByID(tx, tag.ID); err != nil {
		return err
	}
	query := "UPDATE tags SET name=?,alias=? WHERE id=? LIMIT 1"
	_, err := tx.Exec(query, tag.Name, tag.Alias, tag.ID)
	if err != nil {
		return err
	}
	return nil
}

// ListTagsWithCount list all tags with associated post count.
func (tm *TagManager) ListTagsWithCount(tx Querier, limit int64, mergeAlias bool) ([]*TagWithCount, error) {
	query, args := sql_helpers.NewSelect().
		From("post_tags", "pt").From("tags", "t").
		Select("t.*,COUNT(pt.id) size").
		Where("pt.tag_id=t.id").
		GroupBy("t.id").
		OrderBy("size DESC").
		Limit(limit).
		SQL()
	rows, err := tx.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rootTags := make([]*TagWithCount, 0)
	aliasTags := make([]*TagWithCount, 0)
	rootMap := make(map[int64]*TagWithCount, 0)

	for rows.Next() {
		var tag TagWithCount
		err = rows.Scan(&tag.ID, &tag.Name, &tag.Alias, &tag.Count)
		if err != nil {
			return nil, err
		}
		if !mergeAlias {
			rootTags = append(rootTags, &tag)
		} else {
			if tag.Alias == 0 {
				rootTags = append(rootTags, &tag)
				rootMap[tag.ID] = &tag
			} else {
				aliasTags = append(aliasTags, &tag)
			}
		}
	}

	if mergeAlias {
		for _, tag := range aliasTags {
			if root, ok := rootMap[tag.Alias]; ok {
				root.Count += tag.Count
			}
		}
	}

	return rootTags, rows.Err()
}
