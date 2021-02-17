import './styles/styles.scss'
import {BrowserRouter as Router, Route, Routes} from 'react-router-dom';
import Hero from '../hero/Hero'
import NavBar from '../navbar/Navbar'
import AboutUs from '../aboutus/AboutUs';
import ContactUs from '../contactus/ContactUs';
import Footer from '../footer/Footer'
import TermsOfUse from '../footer/TermsOfUse'
import PrivacyPolicy from '../footer/PrivacyPolicy'
import Sentinel from '../sentinel/Sentinel';
import Pricing from '../sentinel/pricing/Pricing';
import Overview from '../sentinel/overview/Overview';
import Documentation from '../sentinel/documentation/Documentation';

const App = () => {
  return (
    <Router>
      <NavBar/>
      <Routes>
        <Route path="/"
          element={<Hero/>}/>
        <Route path="/about-us"
          element={<AboutUs/>}/>
        <Route path="/contact-us"
          element={<ContactUs/>}/>
        <Route path="/terms-of-use"
          element={<TermsOfUse/>}/>
        <Route path="/privacy-policy"
          element={<PrivacyPolicy/>}/>
        <Route path="products/sentinel"
          element={<Sentinel/>}>
          <Route path="/pricing" element={<Pricing/>}/>
          <Route path="/overview" element={<Overview/>}/>
          <Route path="/docs" element={<Documentation/>}/>
          <Route path="/" element={<Overview/>}/>      
        </Route>
      </Routes>
      <Footer/>
    </Router>
  );
}

export default App;

