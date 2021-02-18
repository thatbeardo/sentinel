import { Row, Col } from 'react-bootstrap'
import './styles/styles.scss'
import cost from '../../../assets/cost.png'
import model from '../../../assets/data-model.png'
import integrate from '../../../assets/integrate.png'
import cloud from '../../../assets/cloud.png'

const content = () => {
  return (
    <div className="my-5">
      <h2 className="my-5 text-center">Why Sentinel?</h2>
      <Row className="d-flex justify-content-between">
        
          <Col sm={12} md={5} className="mb-4">
            <img src={model} alt="icon unavailable" className="feature-icon mb-3"/>
            <h6 className="mb-3 feature-title">
                Specialized Data Model
            </h6>
            <p className="feature-description">
              We can handle your complex business requirements with ease. Sentinel's data model is 
              easy to use and flexible. Being truly domain agnostic, we can accomodate your complex 
              business use cases
            </p>
          </Col>
          <Col sm={12} md={5} className="mb-4">
            <img src={cost} alt="icon unavailable" className="feature-icon mb-3"/>
            <h6 className="mb-3 feature-title">
                Pay as you use
            </h6>
            <p className="feature-description">
              No contracts. Generous free-tier. Simple Pricing. Incremental costs as you grow. 
              Tracking costs is easy and updated regularly. Pro tier pricing starts at 5 cents 
              per resource / permission.
            </p>
          </Col>
          <Col sm={12} md={5} className="mb-4">
            <img src={integrate} alt="icon unavailable" className="feature-icon mb-3"/>
            <h6 className="mb-3 feature-title">
              Simple to integrate
            </h6>
            <p className="feature-description">
              Sentinel works great with any of your existing identity provider software. You won't have to
              change your identity provider. You won't have to make major changes to integrate sentinel in
              your product suite.</p>
          </Col>
          <Col sm={12} md={5} className="mb-4">
            <img src={cloud} alt="icon unavailable" className="feature-icon mb-3"/>
            <h6 className="mb-3 feature-title">
              Works in the cloud
            </h6>
            <p className="feature-description">
              We scale by leveraging the latest in cloud technology. Thanks to modern advances in cloud 
              computing, you can trust Sentinel to be a reliable authorization partner.</p>
          </Col>
      </Row>
    </div>
  );
}

export default content;