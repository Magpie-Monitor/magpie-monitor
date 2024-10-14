import {ManagmentServiceApiInstance} from 'api/managment-service';
import './Panel.scss';
import googleLogo from 'assets/google-logo.webp';
import magpieLogo from 'assets/magpie-monitor-icon.png';

const LoginPanel = () => {
    const handleGoogleLogin = (e: React.MouseEvent<HTMLButtonElement>) => {
        e.preventDefault();
        ManagmentServiceApiInstance.login();
    };

    return (
        <div className="login-panel">
            <div className="login-panel__header">Sign in</div>
            <div className="login-panel__subheader">Use your company account to sign in</div>
            <img
                src={magpieLogo}
                alt="Magpie logo"
                className="login-panel__image"
                width="280"
            />
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
            <div className="login-panel__subheader">Can’t access the page?</div>
            <div className="login-panel__subheader-green">Contact your administrator.</div>
        </div>
    );
};
export default LoginPanel;
