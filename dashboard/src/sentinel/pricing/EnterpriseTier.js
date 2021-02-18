import { Button, Card, Col } from 'react-bootstrap'
import group from '../../assets/group.png'
import permission from '../../assets/permission.png'
import headset from '../../assets/headset.png'
import slack from '../../assets/slack.png'

export const EnterpriseTier = () => <Col xs={12} md={12} lg={4}>
<Card className="box-shadow mb-3 pricing-card">
  <Card.Header>
    <h4 className="font-weight-normal">Enterprise</h4>
  </Card.Header>
  <Card.Body>
    <Card.Title className="pricing-card-title">
      <div className="list-unstyled my-4">
          <div className="mb-3 d-flex">
          <img alt="icon unavailable" src={group} className="mr-2 icon" />
            <div className="pt-1 mt-1"> 2M Resources</div>
          </div>
          <div className="mb-3 d-flex">
            <img alt="icon unavailable" src={permission} className="mr-2 icon" />
            <div className="pt-1 mt-1">2M Permissions</div>
          </div>
          <div className="mb-3 d-flex">
            <img alt="icon unavailable" src={headset} className="mr-2 icon" />
            <div className="pt-1 mt-1"> Email & Phone Support</div>
          </div>
          <div className="mb-3 d-flex">
            <img alt="icon unavailable" src={slack} className="mr-2 icon" />
            <div className="pt-1 mt-1">Slack Support</div>
          </div>
      </div>  
    </Card.Title>
    <Card.Text>
      <Button className="outline-btn" block size="lg">
      <a className="learn-link" href="/contact-us">Learn More</a>
      </Button>
    </Card.Text>
    <Button className="pricing-btn" size="lg" block>
      <a className="button-link" href="/contact-us">Contact Us</a>
    </Button>
  </Card.Body>
</Card>
</Col>