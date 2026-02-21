package grpcserver

import (
	"inventory-service/constant"
	"inventory-service/proto/pb"
)

func MapDBStatusToPBStatus(dbStatus string) pb.ReservationStatus {
	switch dbStatus {
	case constant.ReservationStatusPending:
		return pb.ReservationStatus_RESERVATION_STATUS_PENDING
	case constant.ReservationStatusConfirmed:
		return pb.ReservationStatus_RESERVATION_STATUS_CONFIRMED
	case constant.ReservationStatusCancelled:
		return pb.ReservationStatus_RESERVATION_STATUS_CANCELLED
	default:
		return pb.ReservationStatus_RESERVATION_STATUS_UNSPECIFIED
	}
}

func MapPBStatusToDBStatus(pbStatus pb.ReservationStatus) string {
	switch pbStatus {
	case pb.ReservationStatus_RESERVATION_STATUS_PENDING:
		return constant.ReservationStatusPending
	case pb.ReservationStatus_RESERVATION_STATUS_CONFIRMED:
		return constant.ReservationStatusConfirmed
	case pb.ReservationStatus_RESERVATION_STATUS_CANCELLED:
		return constant.ReservationStatusCancelled
	default:
		return constant.ReservationStatusUnspecified
	}
}
