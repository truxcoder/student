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

-- name: UpdateStudent :execrows
UPDATE students
SET name = COALESCE(sqlc.narg('name'), name),
    gender = COALESCE(sqlc.narg('gender'), gender),
    birthday = COALESCE(sqlc.narg('birthday'), birthday),
    score = CASE
                WHEN CAST(sqlc.arg('clear_score') AS INTEGER) = 1 THEN NULL
                ELSE COALESCE(sqlc.narg('score'), score)
            END,
    classroom = COALESCE(sqlc.narg('classroom'), classroom),
    updated_at = CURRENT_TIMESTAMP
WHERE id = :id;

-- name: DeleteStudent :execrows
DELETE FROM students
WHERE id = :id;
