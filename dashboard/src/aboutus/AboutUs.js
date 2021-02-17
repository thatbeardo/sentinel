import {Container} from 'react-bootstrap'
import Team from './Team/Team';
import CoreValues from './CoreValues/CoreValues';

const AboutUs = () => (
  <Container className="pt-5 pb-3">
    <div className="mt-4 pt-4 display-4 d-flex justify-content-center">About Us</div>
    <Team />
    <hr />
    <CoreValues />
  </Container>
);

export default AboutUs;
