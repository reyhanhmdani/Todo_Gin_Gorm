package entity

type Todolist struct {
	ID     int64  `gorm:"primaryKey" json:"id"`
	Title  string `gorm:"type:varchar(300)" json:"title"`
	Status bool   `gorm:"default:false" json:"status"`
}

//func (t Todolist) Read(p []byte) (n int, err error) {
//	//TODO implement me
//	panic("implement me")
//}
