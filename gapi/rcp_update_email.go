package gapi

import (
	"context"
	"errors"
	"time"

	db "github.com/18941325888/sb/db/sqlc"
	"github.com/18941325888/sb/pb"
	"github.com/18941325888/sb/util"
	"github.com/18941325888/sb/val"
	"github.com/18941325888/sb/worker"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) UpdateEmail(ctx context.Context, req *pb.UpdateEmailRequest) (*pb.UpdateEmailResponse, error) {
	authPayload, err := server.authorizeUser(ctx, []string{util.DepositorRole, util.BankerRole})
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	violations := validateUpdateEmailRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	if authPayload.Role != util.BankerRole && authPayload.Username != req.GetUsername() {
		return nil, status.Errorf(codes.PermissionDenied, "cannot update other user's info")
	}

	arg := db.UpdateUserParams{
		Username: req.GetUsername(),
		Email: pgtype.Text{
			String: req.GetEmail(),
			Valid:  true,
		},
		IsEmailVerified: pgtype.Bool{
			Bool:  false,
			Valid: true,
		},
	}

	user, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found")

		}
		return nil, status.Errorf(codes.Internal, "failed to update user: %s", err)

	}

	taskPayload := &worker.PayloadSendVerifyEmail{
		Username: user.Username,
	}
	opts := []asynq.Option{
		asynq.MaxRetry(10),
		asynq.ProcessIn(10 * time.Second),
		asynq.Queue(worker.QueueCritical),
	}
	server.taskDistributor.DistributeTaskSendVerifyEmail(ctx, taskPayload, opts...)

	rsp := &pb.UpdateEmailResponse{
		User: convertUser(user),
	}

	return rsp, nil
}

func validateUpdateEmailRequest(req *pb.UpdateEmailRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if err := val.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}

	return violations
}
