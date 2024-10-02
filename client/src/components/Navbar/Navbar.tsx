import { Link, useLocation } from 'react-router-dom';
import './Navbar.scss';

const NAVBAR_LINKS = [{ to: 'examplePage', label: 'Example page' }];

const Navbar = () => {
  const currentPage = useLocation().pathname;
  return (
    <nav className="navbar">
      <div className="navbar__links">
        {NAVBAR_LINKS.map(({ to, label }) => (
          <Link
            key={to}
            to={to}
            className={`navbar__link ${currentPage.includes(to) ? 'navbar__link--active' : ''}`}
          >
            {label}
          </Link>
        ))}
      </div>
    </nav>
  );
};

export default Navbar;
