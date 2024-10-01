import { navigateToGoogleAuth } from 'api/googleAuth';
import './Panel.scss';
import googleLogo from 'assets/google-logo.webp';

const LoginPanel = () => {
  const handleGoogleLogin = (e: React.MouseEvent<HTMLButtonElement>) => {
    e.preventDefault();
    navigateToGoogleAuth();
  };

  return (
    <div className="login-panel">
      <div className="login-panel__header">Sign in</div>
      <div className="login-panel__subheader">Use your company account to sign in</div>
      <div className="login-panel__body">
        <button className="login-panel__body__button" onClick={handleGoogleLogin}>
          <img
            src={googleLogo}
            alt="Google logo"
            className="login-panel__body__button__logo"
            width="35"
          />
          <span className="login-panel__body__button__text">Sign in with Google </span>
        </button>
      </div>
    </div>
  );
};
export default LoginPanel;
