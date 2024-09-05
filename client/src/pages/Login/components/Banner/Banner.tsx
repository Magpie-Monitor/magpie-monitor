import logo from "@/assets/magpie-monitor-logo.png";
import "./Banner.scss";

const LoginBanner = () => {
  return (
    <div className="login-banner">
      <img src={logo} alt="Relay Beaver" className="login-banner__logo" />
      <h1 className="login-banner__title">Magpie Monitor</h1>
      <p className="login-banner__subtitle">
        Reading logs is for the frogs, let's find insights from them{" "}
      </p>
    </div>
  );
};

export default LoginBanner;
