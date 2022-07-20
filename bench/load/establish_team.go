package load

import (
	"context"
	"errors"
	"math/rand"
	"net"
	"sync"

	"github.com/recruit-tech-tech/RISUCON2022Summer/bench/constant"
	"github.com/recruit-tech-tech/RISUCON2022Summer/bench/fails"
	"github.com/recruit-tech-tech/RISUCON2022Summer/bench/model"
	"github.com/recruit-tech-tech/RISUCON2022Summer/bench/random"
	"github.com/recruit-tech/RISUCON2022Summer/bench/client"
	"golang.org/x/sync/errgroup"
	"golang.org/x/xerrors"
)

type company struct {
	teams []*model.Team
	mu    sync.RWMutex
}

var com = new(company)

func (c *company) addTeam(team *model.Team) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.teams = append(c.teams, team)
}

func (c *company) pickTeam() *model.Team {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.teams[rand.Intn(len(c.teams))]
}

func (c *company) pickUser() *model.User {
	return c.pickTeam().Pick()
}

func (com *company) establishTeam(ctx context.Context, makeConcurrentPost bool) (*model.Team, error) {
	team := random.Team()

	var m = sync.Map{}
	for i := 0; i < 4; i++ {
		eg, childCtx := errgroup.WithContext(ctx)
		for j := 0; j < 4; j++ {
			if makeConcurrentPost && (i*4+j) < constant.ConcurrentPostUserNum {
				eg.Go(func() error {
					for {
						u := com.pickUser()
						_, loaded := m.LoadOrStore(u.ID, struct{}{})
						if !loaded {
							team.Add(u)
							return nil
						}
					}
				})
			} else {
				eg.Go(func() error {
					user, err := createUser(childCtx)
					if err != nil {
						if !errors.Is(err, context.DeadlineExceeded) {
							return err
						}
					}
					team.Add(user)
					return nil
				})
			}
		}

		if err := eg.Wait(); err != nil {
			return nil, err
		}
	}

	com.addTeam(team)

	for i := 0; i < 16; i++ {
		go func() {
			user, err := createUser(ctx)
			if err != nil {
				if err := ctx.Err(); err != nil {
					return
				}
				var nerr net.Error
				if xerrors.As(err, &nerr) && (nerr.Timeout()) {
					fails.Add(nerr)
					return
				}
				if xerrors.Is(err, client.ErrServiceUnavailable) {
					fails.Add(fails.Wrap(err, fails.TrivialErrorCode))
					return
				}
				fails.Add(fails.Wrap(err, fails.ApplicationErrorCode))
				return
			}

			team.Add(user)
		}()
	}

	return team, nil
}

func createUser(ctx context.Context) (*model.User, error) {
	user := random.User()
	c, err := client.New(ctx, client.LoaderType)
	if err != nil {
		return nil, err
	}

	if err := c.PostUser(ctx, user); err != nil {
		return nil, err
	}

	me, err := c.GetMe(ctx)
	if err != nil {
		return nil, err
	}
	user.ID = me.ID
	if !user.IsSame(*me) {
		return nil, errors.New("GET /me: 意図しないユーザーを取得しました")
	}

	if err := c.PostLogout(ctx); err != nil {
		return nil, err
	}

	return user, nil
}
