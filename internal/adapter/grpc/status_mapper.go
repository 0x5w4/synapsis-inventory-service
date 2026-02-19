package grpcadapter

import "inventory-service/proto/pb"

const (
	DBStatusPending     = "PENDING"
	DBStatusConfirmed   = "CONFIRMED"
	DBStatusCancelled   = "CANCELLED"
	DBStatusUnspecified = "UNSPECIFIED"
)

func MapDBStatusToPBStatus(dbStatus string) pb.ReservationStatus {
	switch dbStatus {
	case DBStatusPending:
		return pb.ReservationStatus_RESERVATION_STATUS_PENDING
	case DBStatusConfirmed:
		return pb.ReservationStatus_RESERVATION_STATUS_CONFIRMED
	case DBStatusCancelled:
		return pb.ReservationStatus_RESERVATION_STATUS_CANCELLED
	default:
		return pb.ReservationStatus_RESERVATION_STATUS_UNSPECIFIED
	}
}

func MapPBStatusToDBStatus(pbStatus pb.ReservationStatus) string {
	switch pbStatus {
	case pb.ReservationStatus_RESERVATION_STATUS_PENDING:
		return DBStatusPending
	case pb.ReservationStatus_RESERVATION_STATUS_CONFIRMED:
		return DBStatusConfirmed
	case pb.ReservationStatus_RESERVATION_STATUS_CANCELLED:
		return DBStatusCancelled
	default:
		return DBStatusUnspecified
	}
}
