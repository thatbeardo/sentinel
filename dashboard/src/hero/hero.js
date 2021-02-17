import "./styles/styles.scss"
import background from '../assets/background.mp4'
import { Button } from 'react-bootstrap'

function Hero() {
  return (
    <div>
    <div className="video-background">
      <div className="video-wrap">
        <div id="video">
          <video id="bgvid" autoPlay muted playsInline>
            <source src={background} type="video/mp4"/> 
          </video>
        </div>
      </div>
    </div>

    <div className="caption text-center">
      <h1>REUSABILITY FUELS INNOVATION</h1>
      <Button className="outline-btn mt-4 get-started" >
        <a href="/products/sentinel">Get Started</a>
      </Button>
    </div>
    
    </div>
  );
}

export default Hero