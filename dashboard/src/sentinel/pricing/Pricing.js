import './styles/styles.scss'
import { Row } from 'react-bootstrap'
import { FreeTier } from './FreeTier'
import { ProTier } from './ProTier'
import { EnterpriseTier } from './EnterpriseTier'

const Pricing = () => {
  return (
      <div>
      <div className="pricing-header px-3 py-3 pt-3 md-5 pb-md-4 mx-auto text-center">
        <div className="display-4 d-flex justify-content-center mb-3">Pricing</div>
        <span className="lead mb-5">
          Sign up and get started in minutes. Integrating sentinel with your tech stack is as easy as making an API call. No, literally. It's just a bunch of API calls.
        </span>
        </div>
        <div className="card-deck mb-3 text-center d-flex justify-content-center">
          <Row className="pricing-row">
            <FreeTier/>
            <ProTier/>
            <EnterpriseTier/>
          </Row>
        </div>
      </div>
  )
}

export default Pricing
