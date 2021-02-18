import { Container } from 'react-bootstrap'
import { useEffect } from 'react'
import './styles/styles.scss'

function ContactUs() {
  useEffect(() => {
    // Update the document title using the browser API
    const script = document.createElement('script');
    script.src = 'https://js.hsforms.net/forms/v2.js';
    document.body.appendChild(script);
   
    script.addEventListener('load', () => {
    	if(window.hbspt) {
      	window.hbspt.forms.create({
        	portalId: "9171477",
	        formId: "a0aee583-dfd8-4c03-a399-36223ff8ecaf",
          target: '#hubspotForm'
        })
      }
    });
  });

  return (
  <Container className="contact-us py-5 mt-5">
    <div>
      <div className="pt-4 display-4 d-flex justify-content-center">Contact Us</div>
      <div id="hubspotForm"></div>
    </div>
  </Container>)
}

export default ContactUs;
