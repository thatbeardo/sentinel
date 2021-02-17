import { Button, Container } from 'react-bootstrap'
import './styles/styles.scss'

const ContactUs = () => (
  <Container className="contact-us py-5 mt-5">
    <div className="pb-3 display-4 d-flex justify-content-center">Contact Us</div>
    <form>
      <div className="form-row">
        <div className="form-group col-sm-12 col-md-6">
          <label htmlFor="first-name">First Name</label>
          <input type="email" className="form-control form-input" id="first-name"/>
        </div>
        <div className="form-group col-sm-12 col-md-6">
          <label htmlFor="last-name">Last Name</label>
          <input type="password" className="form-control form-input" id="last-name" />
        </div>
      </div>

      <div className="form-row">
        <div className="form-group col-md-12">
          <label htmlFor="email">Email</label>
          <input type="email" className="form-control form-input" id="email" />
        </div>
      </div>

      <div className="form-row">
        <div className="form-group col-md-12">
          <label htmlFor="company">Company (Optional)</label>
          <input type="email" className="form-control form-input" id="company" />
        </div>
      </div>
        <Button className="pricing-btn mt-3" variant="primary" size="lg" block>Submit</Button>
    </form>
  </Container>
);

export default ContactUs;
