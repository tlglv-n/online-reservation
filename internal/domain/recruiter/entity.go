package recruiter

type Entity struct {
	ID       string  `db:"id" bson:"_id"`
	FullName *string `db:"full_name" bson:"full_name"`
	Email    *string `db:"email" bson:"email"`
	Phone    *int    `db:"phone" bson:"phone"`
}
