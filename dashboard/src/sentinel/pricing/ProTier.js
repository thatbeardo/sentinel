import { Button, Card, Col } from 'react-bootstrap'
import group from '../../assets/group.png'
import permission from '../../assets/permission.png'
import headset from '../../assets/headset.png'
import slack from '../../assets/slack.png'

export const ProTier = () => <Col xs={12} md={12} lg={4}>
<Card className="box-shadow mb-3 pricing-card">
  <Card.Header>
    <h4 className="font-weight-normal">Pro Tier</h4>
  </Card.Header>
  <Card.Body>
    <Card.Title className="pricing-card-title">
      <div className="list-unstyled my-4">
          <div className="mb-3 d-flex">
            <img alt="icon unavailable" src={group} className="mr-2 icon" />
            <div className="pt-1 mt-1">Per 1000 Resources</div>
          </div>
          <div className="mb-3 d-flex">
            <img alt="icon unavailable" src={permission} className="mr-2 icon" />
            <div className="pt-1 mt-1">Per 1000 Permissions</div>
          </div>
          <div className="mb-3 d-flex">
            <img alt="icon unavailable" src={headset} className="mr-2 icon" />
            <div className="pt-1 mt-1">Email & Phone Support</div>
          </div>
          <div className="mb-3 d-flex">
            <img alt="icon unavailable" src={slack} className="mr-2 icon" />
            <div className="pt-1 mt-1">Slack Support</div>
          </div>
      </div>  
    </Card.Title>
    <Card.Text>
        <h2 className="cost">$100<small className="text-muted"> /mo</small></h2>
    </Card.Text>
    <Button className="pricing-btn" variant="primary" size="lg" block>
      <a className="button-link" href="/contact-us">Get Started</a>
    </Button>
  </Card.Body>
</Card>
<small className="pricing-appendix">*Or 5 cents per resource or permission</small>
</Col>