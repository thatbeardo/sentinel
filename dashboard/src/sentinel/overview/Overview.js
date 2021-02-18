import logo from "../../assets/sentinel.png";
import Content from "./content/Content";
import './styles/styles.scss'

const Overview = () => {
  return (
    <div className="text-center">
        <img className="sentinel-logo" src={logo} alt="React logo" />
        <p className="lead mb-5">
          A domain agnostic, application layer authorization solution.
        </p>
        <hr />
      <Content/>
    </div>
  )
}

export default Overview