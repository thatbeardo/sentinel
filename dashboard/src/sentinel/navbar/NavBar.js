import {Navbar, Nav} from 'react-bootstrap'
import brand from '../../assets/helm-brand.png'
import './styles/styles.scss'

const NavBar = () => {

  return (
    <Navbar 
      sticky="top" 
      variant="light" 
      className="sentinel-nav"
      collapseOnSelect  
      expand="lg">
      <Navbar.Brand href="/products/sentinel"><img src={brand}
          alt="Server unavailable"/></Navbar.Brand>
      <Navbar.Toggle className="sentinel-toggle" aria-controls="basic-navbar-nav"/>
      <Navbar.Collapse id="basic-navbar-nav">
        <Nav className="mr-auto ml-4">
          <Nav.Link className="sentinel-nav-link" href="/products/sentinel/overview">Overview</Nav.Link>
          <Nav.Link className="sentinel-nav-link" href="/products/sentinel/pricing">Pricing</Nav.Link>
          <Nav.Link className="sentinel-nav-link" href="/products/sentinel/docs">Documentation</Nav.Link>
        </Nav>
      </Navbar.Collapse>
    </Navbar>
  )
}

export default NavBar