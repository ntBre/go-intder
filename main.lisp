(defparameter bohr 0.529177249d0)
(defparameter ref
  #2A((0.0000000000d0        1.4313902070d0        0.9860411840d0)
      (0.0000000000d0        0.0000000000d0       -0.1242384530d0)
      (0.0000000000d0       -1.4313902070d0        0.9860411840d0)))

(defun row (mat row)
  (let ((ret nil))
    (destructuring-bind (_ cols) (array-dimensions mat)
      (dotimes (c cols)
	(push (aref mat row c) ret)))
    (nreverse ret)))

(defun vec (a b)
  "vec returns the vector from a to b"
  (mapcar #'- b a))

(defun len (vec)
  (labels ((square (x) (* x x)))
    (sqrt (apply #'+ (mapcar #'square vec)))))

(defun sum (lst)
  (apply #'+ lst))

(defun dot (a b)
  (sum (mapcar #'* a b)))

(defun ang (a b)
  "angle between a and b"
  (acos (/ (dot a b)
	   (* (len a)
	      (len b)))))

(defun unit (a)
  (let ((l (len a)))
    (mapcar #'(lambda (x) (if (= 0 l) 0 (/ x l))) a)))

(defun tobohr (a) (* a bohr))
(defun toang (a) (* a (/ 180d0 pi)))

(defparameter h1 (vec (row ref 1) (row ref 0)))
(defparameter h2 (vec (row ref 1) (row ref 2)))

(tobohr (len h2))

(toang (ang h1 h2))

(unit (vec (row ref 0) (row ref 1)))
(unit (vec (row ref 1) (row ref 0)))
(unit (vec (row ref 0) (row ref 0)))

(unit (vec (row ref 1) (row ref 1)))
(unit (vec (row ref 1) (row ref 2)))
(unit (vec (row ref 2) (row ref 1)))

(len (list 0 0 1.64))
