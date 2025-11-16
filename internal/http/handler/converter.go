package handler

import "github.com/dinoagera/AvitoPullRequest/internal/domain"

func dtoToDomainTeam(req AddTeamRequest) *domain.Team {
	team := &domain.Team{
		Name:    req.TeamName,
		Members: make([]domain.User, len(req.Members)),
	}
	for i, m := range req.Members {
		team.Members[i] = domain.User{
			ID:       m.UserID,
			Name:     m.Username,
			TeamName: req.TeamName,
			IsActive: m.IsActive,
		}
	}
	return team
}
func domainTodtoTeam(team *domain.Team) interface{} {
	members := make([]map[string]interface{}, len(team.Members))
	for i, u := range team.Members {
		members[i] = map[string]interface{}{
			"user_id":   u.ID,
			"username":  u.Name,
			"is_active": u.IsActive,
		}
	}
	return map[string]interface{}{
		"team_name": team.Name,
		"members":   members,
	}
}
func domainTodtoUser(user *domain.User) interface{} {
	return map[string]interface{}{
		"user_id":   user.ID,
		"username":  user.Name,
		"team_name": user.TeamName,
		"is_active": user.IsActive,
	}
}
func dtoToDomainPR(req CreateRPRequest) *domain.PullRequest {
	rp := &domain.PullRequest{
		ID:       req.ID,
		Name:     req.Name,
		AuthorID: req.AuthorID,
	}
	return rp
}
