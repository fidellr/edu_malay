package repository

import (
	"context"
	"strings"

	"github.com/fidellr/edu_malay/utils"
	"github.com/globalsign/mgo/bson"

	"github.com/fidellr/edu_malay/model"

	"github.com/fidellr/edu_malay/model/clc"
	"github.com/globalsign/mgo"
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

func (r *ClcMongo) Remove(ctx context.Context, id string) error {
	sess := r.Session.Clone()
	defer sess.Close()

	idBson := bson.ObjectIdHex(id)
	return sess.DB(r.DBName).C(clcCollectionName).Remove(bson.M{"_id": idBson})
}
