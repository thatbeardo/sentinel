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
            <p className="feature-description">We can handle your complex business requirements with ease.</p>
          </Col>
          <Col sm={12} md={5} className="mb-4">
            <img src={cost} alt="icon unavailable" className="feature-icon mb-3"/>
            <h6 className="mb-3 feature-title">
                Pay as you use
            </h6>
            <p className="feature-description">No contracts. Generous free-tier. Simple Pricing. Incremental costs as you grow.</p>
          </Col>
          <Col sm={12} md={5} className="mb-4">
            <img src={integrate} alt="icon unavailable" className="feature-icon mb-3"/>
            <h6 className="mb-3 feature-title">
              Simple to integrate
            </h6>
            <p className="feature-description">Sentinel works great with any of your existing identity provider software.</p>
          </Col>
          <Col sm={12} md={5} className="mb-4">
            <img src={cloud} alt="icon unavailable" className="feature-icon mb-3"/>
            <h6 className="mb-3 feature-title">
              Works in the cloud
            </h6>
            <p className="feature-description">We scale by leveraging the latest in cloud technology.</p>
          </Col>
      </Row>
    </div>
  );
}

export default content;