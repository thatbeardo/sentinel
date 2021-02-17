import {Container, Row} from 'react-bootstrap'
import './styles/styles.scss'

const email = 'contact@bithippie.com'
const subject = 'Sentinel Inquiry'
const body = 'I wish to learn more about Sentinel and it\'s abilities'

const Footer = () => (
  <footer className="bg-light footer py-4">
    <Container>
      <Row className="d-flex justify-content-between">
        <div>
          <a href="/privacy-policy">Privacy Policy
          </a>
          &nbsp;&nbsp;&nbsp;| &nbsp;&nbsp;&nbsp;
          <a href="/terms-of-use">Terms of Use</a>
        </div>
        <div>
          Say hello: &nbsp;
          <a href={
            `mailto:${email}?subject=${
              encodeURIComponent(subject) || ""
            }&body=${
              encodeURIComponent(body) || ""
            }`
          }>contact@bithippe.com</a>
        </div>
      </Row>

      <Row>
        &copy; {
        new Date().getFullYear()
      }
        Copyright BitHippie. All Rights
                  Reserved.
      </Row>

    </Container>
  </footer>
);

export default Footer;

