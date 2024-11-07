import { Link, useLocation } from 'react-router-dom';
import './NavbarTab.scss';
import SVGIcon from 'components/SVGIcon/SVGIcon';

export interface NavbarTabProps {
  label: string;
  destination: string;
  iconName: string;
  onClick?: () => void;
}

const NavbarTab = ({
  label,
  destination,
  iconName,
  onClick = () => {},
}: NavbarTabProps) => {
  const location = useLocation();

  return (
    <Link
      className={`navbar-tab ${location.pathname == destination ? 'navbar-tab--selected' : ''}`}
      to={destination}
      onClick={onClick}
    >
      <SVGIcon iconName={iconName} />
      <div className="navbar-tab__link-wrapper">
        <div key={destination} className={'navbar-tab__link'}>
          {label}
        </div>
      </div>
    </Link>
  );
};

export default NavbarTab;
