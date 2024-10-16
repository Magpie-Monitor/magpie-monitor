import './Navbar.scss';
import magpieMonitorLogo from 'assets/magpie-monitor-icon.png';
import NavbarTab from './NavbarTab/NavbarTab';
import { ManagmentServiceApiInstance } from 'api/managment-service';

const Navbar = () => {
  const signOut = () => {
    ManagmentServiceApiInstance.logout();
  };

  return (
    <nav className="navbar">
      <div>
        <div className="navbar__logo">
          <img src={magpieMonitorLogo} alt="Magpie Monitor logo" className="navbar__logo__image" />
          <div className="navbar__logo__name">
            <div>Magpie</div>
            <div>Monitor</div>
          </div>
        </div>
        <div className="navbar__links">
          <NavbarTab label={'Dashboard'} destination={'/dashboard'} iconName={'dashboard-icon'} />
          <NavbarTab label={'Reports'} destination={'/reports'} iconName={'reports-icon'} />
          <NavbarTab label={'Settings'} destination={'/settings'} iconName={'setting-icon'} />
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
