package repository

import (
	"context"
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
	var teacherM teacher.ProfileEntity
	var clcM clc.ProfileEntity
	teacherIdntity := new(assembler.TeacherIdentity)
	clcIDBson := bson.ObjectIdHex(clcID)
	done := make(chan bool)

	go func(teacherId bson.ObjectId) {
		err = sess.DB(r.DBName).C(teacherCollection).Find(bson.M{"_id": teacherId}).One(&teacherM)
		if err != nil {
			done <- true
			return
		}
		done <- true
	}(m.TeacherID)

	go func(clcID string) {
		err = sess.DB(r.DBName).C(clcCollection).Find(bson.M{"_id": clcIDBson}).One(&clcM)
		if err != nil {
			done <- true
			return
		}

		done <- true
	}(clcID)

	if d, _ := <-done; d {
		if err != nil {
			return err
		}

		teacherIdntity.ID = teacherM.ID
		teacherIdntity.FirstName = teacherM.FirstName
		teacherIdntity.LastName = teacherM.LastName
		teacherIdntity.Gender = teacherM.Gender
		teacherIdntity.StartWorkDate = m.StartWorkDate
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
