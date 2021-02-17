import {Outlet} from 'react-router-dom';
import NavBar from './navbar/NavBar'
import { Container } from 'react-bootstrap'

const Sentinel = () => {
  return (
    <div className="pt-5 mt-3">
      <NavBar />
      <Container className="my-5 pb-2">
        <Outlet/>
      </Container>
    </div>
  )
}

export default Sentinel
