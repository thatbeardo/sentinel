import data from './data'
import { Row, Col } from 'react-bootstrap'
import './styles/styles.scss'

const content = () => {
  return (
    <div className="my-5">
      <h2 className="my-5 text-center">Why Sentinel?</h2>
      <Row className="d-flex justify-content-between">
        {data.map((col, i) => (
          <Col key={i} md={5} className="mb-4">
            <h6 className="mb-3">
                {col.title}
            </h6>
            <p>{col.description}</p>
          </Col>
        ))}
      </Row>
    </div>
  );
}

export default content;