package ford

import "context"

type Provisioner struct {
}

func NewProvisioner() *Provisioner {
	return &Provisioner{}
}

func (p *Provisioner) Provision(ctx context.Context, pipeline *Pipeline) error {
	for _, team := range pipeline.Teams {
		err := p.provisionTeam(ctx, pipeline, team)

		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Provisioner) provisionTeam(ctx context.Context, pipeline *Pipeline, team Team) error {

}
