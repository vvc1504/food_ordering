package spec

type CouponValidatorStatusCode StatusCode

const (
	CouponValidatorNone    CouponValidatorStatusCode = 0
	CouponValidatorInvalid CouponValidatorStatusCode = 1
)

type ErrorCouponValidatorStatus Status

// CouponValidator defines the logic for validating promo codes.
type CouponValidator interface {
	Validate(Code string) (Resp bool, Err ErrorCouponValidatorStatus)
}
