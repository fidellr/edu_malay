package repository

import (
	"context"

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

func (r *TeacherMongo) FindAll(ctx context.Context, filter *teacher.Filter) ([]*teacher.ProfileEntity, string, error) {
	// var teachers []*teacher.ProfileEntity
	teachers := make([]*teacher.ProfileEntity, 0)
	query := make(bson.M)
	sess := r.Session.Clone()
	defer sess.Close()

	if filter.Cursor != "" {
		createdAt, err := utils.DecodeTime(filter.Cursor)
		if err != nil {
			return teachers, "", err
		}
		query["created_at"] = bson.M{"$lt": createdAt}
	}

	err := sess.DB(r.DBName).C(teacherCollectionName).Find(query).Limit(filter.Num).Sort("-created_at").All(&teachers)
	if err != nil {
		return teachers, "", err
	}

	if len(teachers) < 1 {
		return teachers, "", nil
	}

	lastUsers := teachers[len(teachers)-1]
	nextCursors := utils.EncodeTime(lastUsers.CreatedAt)

	return teachers, nextCursors, nil
}

func (r *TeacherMongo) GetByID(ctx context.Context, id string) (*teacher.ProfileEntity, error) {
	t := new(teacher.ProfileEntity)
	sess := r.Session.Clone()
	defer sess.Close()

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
