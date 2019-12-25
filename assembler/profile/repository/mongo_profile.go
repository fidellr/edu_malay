package repository

import (
	"context"
	"sync"
	"time"

	"github.com/fidellr/edu_malay/model/clc"

	"github.com/fidellr/edu_malay/model/teacher"

	"github.com/globalsign/mgo/bson"

	"github.com/fidellr/edu_malay/model/assembler"

	"github.com/globalsign/mgo"
)

type ProfileAssemblerMongo struct {
	Session *mgo.Session
	DBName  string
}

var (
	profileAssemblerCollection = "profile-assembler"
	teacherCollection          = "teacher"
	clcCollection              = "clc"
)

func NewProfileAssemblerMongo(session *mgo.Session, DBName string) *ProfileAssemblerMongo {
	return &ProfileAssemblerMongo{
		session,
		DBName,
	}
}

func (r *ProfileAssemblerMongo) Create(ctx context.Context, clcID string, m *assembler.ProfileAssemblerParam) error {
	sess := r.Session.Clone()
	defer sess.Close()

	var err error
	var w sync.WaitGroup
	done := make(chan bool)
	w.Add(2)

	teacherIdntity := new(assembler.TeacherIdentity)
	clcIDBson := bson.ObjectIdHex(clcID)

	go func(teacherId bson.ObjectId, co chan<- bool) {
		teacherM := new(teacher.ProfileEntity)
		if teacherM.ID == "" {
			err = sess.DB(r.DBName).C(teacherCollection).Find(bson.M{"_id": teacherId}).One(teacherM)
			if err != nil {
				w.Done()
				co <- true
				return
			}

			teacherIdntity.ID = teacherM.ID
			teacherIdntity.FirstName = teacherM.FirstName
			teacherIdntity.LastName = teacherM.LastName
			teacherIdntity.Gender = teacherM.Gender
			teacherIdntity.StartWorkDate = m.StartWorkDate
		}
		w.Done()
		co <- true
	}(m.TeacherID, done)

	go func(co chan<- bool) {
		clcM := new(clc.ProfileEntity)
		if clcM.ID == "" {
			err = sess.DB(r.DBName).C(clcCollection).Find(bson.M{"_id": clcIDBson}).One(clcM)
			if err != nil {
				w.Done()
				co <- true
				return
			}
		}

		w.Done()
		co <- true
	}(done)

	w.Wait()
	if d := <-done; d {
		if err != nil {
			return err
		}

		err = sess.DB(r.DBName).C(profileAssemblerCollection).Update(bson.M{"clc_id": clcIDBson}, bson.M{"$addToSet": bson.M{"teachers": teacherIdntity}})
		if err == mgo.ErrNotFound {
			assemblerEnt := new(assembler.ProfileAssemblerEntity)

			assemblerEnt.ClcID = clcIDBson
			assemblerEnt.CreatedAt = time.Now()
			assemblerEnt.UpdatedAt = time.Now()
			assemblerEnt.Teachers = append(assemblerEnt.Teachers, *teacherIdntity)

			err = sess.DB(r.DBName).C(profileAssemblerCollection).Insert(assemblerEnt)
			if err != nil {
				return err
			}
		}
	}

	return err
}

func (r *ProfileAssemblerMongo) FetchAll(ctx context.Context) ([]*assembler.ProfileAssemblerEntity, error) {
	sess := r.Session.Clone()
	defer sess.Close()
	assembledProfileEnt := make([]*assembler.ProfileAssemblerEntity, 0)

	err := sess.DB(r.DBName).C(profileAssemblerCollection).Find(bson.M{}).All(&assembledProfileEnt)
	if err != nil {
		return assembledProfileEnt, err
	}

	return assembledProfileEnt, err
}

func (r *ProfileAssemblerMongo) GetByID(ctx context.Context, id string) (*assembler.ProfileAssemblerEntity, error) {
	sess := r.Session.Clone()
	defer sess.Close()
	assembledProfileEnt := new(assembler.ProfileAssemblerEntity)

	idBson := bson.ObjectIdHex(id)
	err := sess.DB(r.DBName).C(profileAssemblerCollection).Find(bson.M{"_id": idBson}).One(assembledProfileEnt)
	if err != nil {
		return assembledProfileEnt, err
	}

	return assembledProfileEnt, nil
}

func (r *ProfileAssemblerMongo) Update(ctx context.Context, id string, teacherParam *assembler.ProfileAssemblerParam, isEditing bool) error {
	sess := r.Session.Clone()
	defer sess.Close()
	idBson := bson.ObjectIdHex(id)

	if !isEditing {
		err := sess.DB(r.DBName).C(profileAssemblerCollection).Update(bson.M{"_id": idBson}, bson.M{"$pull": bson.M{"teachers": bson.M{"_id": teacherParam.TeacherID}}})
		if err != nil {
			return err
		}
		return nil
	}

	err := sess.DB(r.DBName).C(profileAssemblerCollection).Update(bson.M{"_id": idBson, "teachers._id": teacherParam.TeacherID}, bson.M{"$set": bson.M{"teachers.$.start_work_date": teacherParam.StartWorkDate}})
	if err != nil {
		return err
	}

	return nil
}

func (r *ProfileAssemblerMongo) Remove(ctx context.Context, assmblrProfileID string) error {
	sess := r.Session.Clone()
	defer sess.Close()

	idBson := bson.ObjectIdHex(assmblrProfileID)
	err := sess.DB(r.DBName).C(profileAssemblerCollection).Remove(bson.M{"_id": idBson})
	if err != nil {

		return err
	}

	return nil
}
