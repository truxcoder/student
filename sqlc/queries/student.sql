-- name: GetStudentByID :one
SELECT id, name, gender, birthday, score, classroom, created_at, updated_at
FROM students
WHERE id = :id;

-- name: ListStudents :many
SELECT id, name, gender, birthday, score, classroom, created_at, updated_at
FROM students;

-- name: CreateStudent :exec
INSERT INTO students (id, name, gender, birthday, score, classroom)
VALUES (:id, :name, :gender, :birthday, :score, :classroom);

-- name: UpdateStudent :exec
UPDATE students
SET name = :name,
    gender = :gender,
    birthday = :birthday,
    score = COALESCE(sqlc.narg('score'), score),
    classroom = :classroom,
    updated_at = CURRENT_TIMESTAMP
WHERE id = :id;

-- name: DeleteStudent :exec
DELETE FROM students
WHERE id = :id;
