import './Navbar.scss';
import magpieMonitorLogo from 'assets/magpie-monitor-icon.png';
import NavbarTab from './NavbarTab/NavbarTab';
import { ManagmentServiceApiInstance } from 'api/managment-service';
import { useNavigate } from 'react-router-dom';
import { useAuth } from 'providers/AuthProvider/AuthProvider.tsx';

const Navbar = () => {
  const navigate = useNavigate();
  const { setAuthenticationInfo } = useAuth();
  const signOut = () => {
    ManagmentServiceApiInstance.logout();
    setAuthenticationInfo(null);
  };

  const handleHomeNavigation = () => {
    navigate('/');
  };

  return (
    <nav className="navbar">
      <div>
        <div className="navbar__logo">
          <img
            src={magpieMonitorLogo}
            alt="Magpie Monitor logo"
            className="navbar__logo__image"
            onClick={handleHomeNavigation}
          />
          <div className="navbar__logo__name">
            <div>Magpie</div>
            <div>Monitor</div>
          </div>
        </div>
        <div className="navbar__links">
          <NavbarTab
            label={'Dashboard'}
            destination={'/'}
            iconName={'dashboard-icon'}
          />
          <NavbarTab
            label={'Reports'}
            destination={'/reports'}
            iconName={'reports-icon'}
          />
          <NavbarTab
            label={'Clusters'}
            destination={'/clusters'}
            iconName={'clusters-icon--white'}
          />
          <NavbarTab
            label={'Notifications'}
            destination={'/notifications'}
            iconName={'notification-icon--white'}
          />
        </div>
      </div>
      <div className="navbar__sign-out ">
        <NavbarTab
          label={'Sign Out'}
          destination={'/login'}
          iconName={'sign-out-icon'}
          onClick={signOut}
        />
      </div>
    </nav>
  );
};

export default Navbar;
