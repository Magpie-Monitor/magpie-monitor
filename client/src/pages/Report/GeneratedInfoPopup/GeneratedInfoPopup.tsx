import React from 'react';
import { useNavigate } from 'react-router-dom';
import OverlayComponent from 'components/OverlayComponent/OverlayComponent.tsx';
import './GeneratedInfoPopup.scss';
import ActionButton, {ActionButtonColor} from 'components/ActionButton/ActionButton.tsx';

interface GeneratedInfoPopupProps {
    isDisplayed: boolean;
    onClose: () => void;
}

const GeneratedInfoPopup: React.FC<GeneratedInfoPopupProps> = ({ isDisplayed, onClose }) => {
    const navigate = useNavigate();

    const handleClose = () => {
        onClose();
        navigate('/dashboard');
    };

    return (
        <OverlayComponent isDisplayed={isDisplayed} onClose={onClose}>
            <div className="report-generated-info-popup">
                <h2>Report is being generated</h2>
                {/* eslint-disable-next-line max-len */}
                <p>Your report for the selected dates is being generated and will be available within 24 hours.</p>
                <div className="report-generated-info-popup__button">
                    <ActionButton onClick={handleClose}
                                  description="Close" color={ActionButtonColor.RED}/>
                </div>
            </div>
        </OverlayComponent>
    );
};

export default GeneratedInfoPopup;
