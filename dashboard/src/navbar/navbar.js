import './styles/styles.scss'
import React, { useState } from 'react'
import {Navbar, Nav, NavDropdown} from 'react-bootstrap'
import brand from '../assets/brand.png'
import { useLocation } from 'react-router-dom';

function NavBar() {

  let location = useLocation();

  const [navbarTransparency, setNavbarTransparency] = useState(false)
  const updateTransparency = () => setNavbarTransparency(window.scrollY >= 20)
  
  const getNavbarBackgroundStyle = () => { 
    var style = ''
    style += location.pathname !== '/' ? 'active' : ''
    style += navbarTransparency ? ' navbar-hide' : ''
    return style
  }

  window.addEventListener('scroll', updateTransparency)

  return (

    <Navbar className={
        getNavbarBackgroundStyle()
      }
      collapseOnSelect  
      expand="lg"
      fixed="top"
      variant="dark">
      <Navbar.Brand href="/"><img src={brand}
          alt="Server unavailable"/></Navbar.Brand>
      <Navbar.Toggle className="bithippie-toggle" aria-controls="basic-navbar-nav"/>
      <Navbar.Collapse id="basic-navbar-nav">
        <Nav className="mr-auto">
          <NavDropdown title="Products" id="basic-nav-dropdown">
            <NavDropdown.Item href="/products/sentinel">Sentinel</NavDropdown.Item>
          </NavDropdown>
          <Nav.Link href="/about-us">About Us</Nav.Link>
          <Nav.Link href="/contact-us">Contact Us</Nav.Link>
        </Nav>
      </Navbar.Collapse>
    </Navbar>
  );
}

export default NavBar
