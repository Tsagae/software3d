package basics

// Order of transformation = Scale, rotate, translate
type Transform struct {
	Scaling     Scalar
	Rotation    Quaternion
	Translation Vector3
}

// Constructors
func NewTransform(scaling Scalar, rotation Quaternion, translation Vector3) Transform {
	return Transform{Scaling: scaling, Rotation: rotation, Translation: translation}
}

func NewZeroTransform() Transform {
	return Transform{Scaling: 1, Rotation: NewIdentityQuaternion(), Translation: Vector3{}}
}

/*
func NewForwardTransform() *Transform {
	return &Transform{Scale: 1, Rotation: *NewForwardQuaternion(), Position: *ZeroVector(), inverse: false}
}

func NewBackwardTransform() *Transform {
	return &Transform{Scale: 1, Rotation: *NewBackWardQuaternion(), Position: *ZeroVector(), inverse: false}
}
*/

// Operations that modify this
func (t *Transform) ThisCumulate(t2 *Transform) {
	t.Scaling *= t2.Scaling
	//b.Rotation = b.Rotation.Mul(&t2.Rotation)
	t.Rotation = t2.Rotation.Mul(&t.Rotation)
	t.Translation = t2.Rotation.Rotated(t.Translation.Mul(t2.Scaling)).Add(&t2.Translation)
}

func (t *Transform) ThisInvert() {
	t.Scaling = 1 / t.Scaling
	t.Rotation.ThisConjugate()
	t.Translation = t.Rotation.Rotated(t.Translation.Inverse().Mul(t.Scaling))
}

// Operations that do not change this
func (t *Transform) Cumulate(t2 *Transform) Transform {
	tNew := *t
	tNew.ThisCumulate(t2)
	return tNew
}

func (t *Transform) Inverse() Transform {
	copy := *t
	copy.ThisInvert()
	return copy
}

// modifies v
func (t *Transform) ApplyToVector(v *Vector3) {
	//return r.apply_to( s * v ); // no traslation
	*v = t.Rotation.Rotated(v.Mul(t.Scaling))
}

// modifies p
func (t *Transform) ApplyToPoint(p *Vector3) {
	//return r.apply_to( s * p ) + t;
	*p = t.Rotation.Rotated(p.Mul(t.Scaling)).Add(&t.Translation)
}

func (t *Transform) Equals(t2 *Transform) bool {
	return t.Scaling.Equals(t2.Scaling) && t.Rotation.Equals(&t2.Rotation) && t.Translation.Equals(&t2.Translation)
}
