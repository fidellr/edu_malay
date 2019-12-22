package repository

import (
	"context"
	"strings"

	"github.com/fidellr/edu_malay/model"
	"github.com/fidellr/edu_malay/model/teacher"
	"github.com/fidellr/edu_malay/utils"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type TeacherMongo struct {
	Session *mgo.Session
	DBName  string
}

var (
	teacherCollectionName = "teacher"
)

func NewTeacherMongo(session *mgo.Session, DBName string) *TeacherMongo {
	return &TeacherMongo{
		session,
		DBName,
	}
}

func (r *TeacherMongo) Create(ctx context.Context, t *teacher.ProfileEntity) error {
	sess := r.Session.Clone()
	defer sess.Close()

	err := sess.DB(r.DBName).C(teacherCollectionName).Insert(t)
	if err != nil {
		return err
	}

	return nil
}

func (r *TeacherMongo) FindAll(ctx context.Context, filter *model.Filter) ([]*teacher.ProfileEntity, string, error) {
	sess := r.Session.Clone()
	defer sess.Close()

	teachers := make([]*teacher.ProfileEntity, 0)
	query := make(bson.M)

	if filter.Cursor != "" {
		createdAt, err := utils.DecodeTime(filter.Cursor)
		if err != nil {
			return teachers, "", err
		}
		query["created_at"] = bson.M{"$lt": createdAt}
	}

	col := sess.DB(r.DBName).C(teacherCollectionName)

	if filter.Search != "" {
		query["$text"] = bson.M{"$search": strings.ToLower(filter.Search)}
	}

	q := col.Find(query).Sort("-created-at")
	if filter.Search != "" {
		q = q.Select(bson.M{"score": bson.M{"$meta": "textScore"}})
	}

	err := q.Limit(filter.Num).All(&teachers)
	if err != nil {
		return teachers, "", err
	}

	if len(teachers) < 1 {
		return teachers, "", nil
	}

	lastData := teachers[len(teachers)-1]
	nextCursors := utils.EncodeTime(lastData.CreatedAt)
	return teachers, nextCursors, nil
}

func (r *TeacherMongo) GetByID(ctx context.Context, id string) (*teacher.ProfileEntity, error) {
	sess := r.Session.Clone()
	defer sess.Close()
	t := new(teacher.ProfileEntity)

	idBson := bson.ObjectIdHex(id)
	err := sess.DB(r.DBName).C(teacherCollectionName).Find(bson.M{"_id": idBson}).One(t)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (r *TeacherMongo) Update(ctx context.Context, id string, t *teacher.ProfileEntity) error {
	sess := r.Session.Clone()
	defer sess.Close()

	idBson := bson.ObjectIdHex(id)
	return sess.DB(r.DBName).C(teacherCollectionName).Update(bson.M{"_id": idBson}, t)
}

func (r *TeacherMongo) Remove(ctx context.Context, id string) error {
	sess := r.Session.Clone()
	defer sess.Close()

	idBson := bson.ObjectIdHex(id)
	return sess.DB(r.DBName).C(teacherCollectionName).Remove(bson.M{"_id": idBson})
}
