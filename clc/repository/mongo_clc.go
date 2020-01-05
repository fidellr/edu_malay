package repository

import (
	"context"
	"strings"
	"time"

	"github.com/fidellr/edu_malay/model/teacher"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"

	"github.com/fidellr/edu_malay/model"

	"github.com/fidellr/edu_malay/model/clc"
	"github.com/fidellr/edu_malay/utils"
)

type ClcMongo struct {
	Session *mgo.Session
	DBName  string
}

var (
	clcCollectionName = "clc"
)

func NewClcProfileMongo(session *mgo.Session, DBName string) *ClcMongo {
	return &ClcMongo{
		session,
		DBName,
	}
}

func (r *ClcMongo) Create(ctx context.Context, clc *clc.ProfileEntity) error {
	sess := r.Session.Clone()
	defer sess.Close()

	err := sess.DB(r.DBName).C(clcCollectionName).Insert(clc)
	if err != nil {
		return err
	}

	return nil
}

func (r *ClcMongo) FindAll(ctx context.Context, filter *model.Filter) ([]*clc.ProfileEntity, string, error) {
	sess := r.Session.Clone()
	defer sess.Close()

	clcs := make([]*clc.ProfileEntity, 0)
	query := make(bson.M)

	if filter.Cursor != "" {
		createdAt, err := utils.DecodeTime(filter.Cursor)
		if err != nil {
			return clcs, "", err
		}
		query["created_at"] = bson.M{"$lt": createdAt}
	}

	col := sess.DB(r.DBName).C(clcCollectionName)
	if filter.Search != "" {
		query["$text"] = bson.M{"$search": strings.ToLower(filter.Search)}
	}

	q := col.Find(query).Sort("-created-at")
	if filter.Search != "" {
		q = q.Select(bson.M{"score": bson.M{"$meta": "textScore"}})
	}

	err := q.Limit(filter.Num).All(&clcs)
	if err != nil {
		return clcs, "", err
	}

	if len(clcs) < 1 {
		return clcs, "", nil
	}

	lastData := clcs[len(clcs)-1]
	nextCursors := utils.EncodeTime(lastData.CreatedAt)

	return clcs, nextCursors, nil
}

func (r *ClcMongo) GetByID(ctx context.Context, id string) (*clc.ProfileEntity, error) {
	sess := r.Session.Clone()
	defer sess.Close()
	clc := new(clc.ProfileEntity)

	idBson := bson.ObjectIdHex(id)
	err := sess.DB(r.DBName).C(clcCollectionName).Find(bson.M{"_id": idBson}).One(clc)
	if err != nil {
		return nil, err
	}

	return clc, nil
}

func (r *ClcMongo) Update(ctx context.Context, id string, clc *clc.ProfileEntity) error {
	sess := r.Session.Clone()
	defer sess.Close()

	idBson := bson.ObjectIdHex(id)
	return sess.DB(r.DBName).C(clcCollectionName).Update(bson.M{"_id": idBson}, clc)
}

func (r *ClcMongo) AssembleProfile(ctx context.Context, clcID, teacherID, startDate string) error {
	sess := r.Session.Clone()
	defer sess.Close()

	clcIDBson := bson.ObjectIdHex(clcID)
	teacherIDBson := bson.ObjectIdHex(teacherID)

	var mT *teacher.ProfileEntity
	err := sess.DB(r.DBName).C("teacher").Find(bson.M{"_id": teacherIDBson}).One(&mT)
	if err == mgo.ErrNotFound {
		return err
	}

	err = sess.DB(r.DBName).C("teacher").Update(bson.M{"_id": teacherIDBson}, bson.M{"$set": bson.M{"updated_at": time.Now(), "start_work_date": startDate, "is_assembled": true}})
	if err != nil {
		return err
	}

	mT.StartWorkDate = startDate
	err = sess.DB(r.DBName).C(clcCollectionName).Update(bson.M{"_id": clcIDBson}, bson.M{"$addToSet": bson.M{"teachers": mT}})
	if err != nil {
		return err
	}

	return nil
}

func (r *ClcMongo) UpdateAssembledProfile(ctx context.Context, clcID, teacherID, startWorkDate string, isEditing bool) error {
	sess := r.Session.Copy()
	defer sess.Close()
	query := make(bson.M)
	selectorQuery := make(bson.M)

	clcIDBson := bson.ObjectIdHex(clcID)
	teacherIDBson := bson.ObjectIdHex(teacherID)
	if !isEditing {
		selectorQuery["_id"] = clcIDBson
		query["$pull"] = bson.M{"teachers": bson.M{"_id": teacherIDBson}}
		err := sess.DB(r.DBName).C("teacher").Update(bson.M{"_id": teacherIDBson}, bson.M{"$set": bson.M{"updated_at": time.Now(), "is_assembled": false, "start_work_date": ""}})
		if err != nil {
			return err
		}
	} else {
		selectorQuery = bson.M{"_id": clcIDBson, "teachers._id": teacherIDBson}
		query["$set"] = bson.M{"teachers.$.start_work_date": startWorkDate}
		err := sess.DB(r.DBName).C("teacher").Update(bson.M{"_id": teacherIDBson}, bson.M{"$set": bson.M{"updated_at": time.Now(), "start_work_date": startWorkDate}})
		if err != nil {
			return err
		}
	}

	return sess.DB(r.DBName).C(clcCollectionName).Update(selectorQuery, query)
}

func (r *ClcMongo) Remove(ctx context.Context, id string) error {
	sess := r.Session.Clone()
	defer sess.Close()

	idBson := bson.ObjectIdHex(id)
	return sess.DB(r.DBName).C(clcCollectionName).Remove(bson.M{"_id": idBson})
}
