package entity

type User struct {
	Id                 string `json:"id,omitempty"`
	Name               string `json:"name" binding:"required"`
	Email              string `json:"email" binding:"required"`
	DateOfBirth        string `json:"date_of_birth" binding:"required"`
	Gender             string `json:"gender" binding:"required"`
	LastAttendanceTime string `json:"last_attendance_time"`
	FirstAttendance    string `json:"first_attendance"`
	TotalAttendance    int    `json:"total_attendance"`
	Password           string `json:"password" binding:"required"`
}

type Token struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type DateOfIMages struct {
	Date string `json:"date" binding:"required"`
}
