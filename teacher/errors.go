package teacher

type ErrorTeacherNotFound struct {
	message string
}

func (err ErrorTeacherNotFound) Error() string {
	return err.message
}

func NewErrorTeacherNotFound(message string) ErrorTeacherNotFound {
	return ErrorTeacherNotFound{message: message}
}

type ErrorInvalidTeacherData struct {
	message string
}

func (err ErrorInvalidTeacherData) Error() string {
	return err.message
}

func NewErrorInvalidTeacherData(message string) ErrorInvalidTeacherData {
	return ErrorInvalidTeacherData{message: message}
}

type ErrorTeacherConflict struct {
	message string
}

func NewErrorConflictUser(message string) ErrorTeacherConflict {
	return ErrorTeacherConflict{
		message: message,
	}
}
