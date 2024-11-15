import React from 'react';
import { useNavigate } from 'react-router-dom';
import OverlayComponent from 'components/OverlayComponent/OverlayComponent.tsx';
import './GeneratedInfoPopup.scss';
import ActionButton, {
    ActionButtonColor,
} from 'components/ActionButton/ActionButton.tsx';
import Hourglass from 'components/Hourglass/Hourglass.tsx';
import {
    REPORT_IN_PROGRESS_TITLE,
    REPORT_IN_PROGRESS_DESCRIPTION,
} from 'messages/info-messages';

interface GeneratedInfoPopupProps {
    isDisplayed: boolean;
    onClose: () => void;
}

const GeneratedInfoPopup: React.FC<GeneratedInfoPopupProps> = ({
    isDisplayed,
    onClose,
}) => {
    const navigate = useNavigate();

    const handleClose = () => {
        onClose();
        navigate('/');
    };

    return (
        <OverlayComponent isDisplayed={isDisplayed} onClose={handleClose}>
            <div className="report-generated-info-popup">
                <div className="report-generated-info-popup__header">
                    <h2 className="report-generated-info-popup__title">
                        {REPORT_IN_PROGRESS_TITLE}
                    </h2>
                    <Hourglass />
                </div>
                <p className="report-generated-info-popup__description">
                    {REPORT_IN_PROGRESS_DESCRIPTION}
                </p>
                <div className="report-generated-info-popup__button">
                    <ActionButton
                        onClick={handleClose}
                        description="Close"
                        color={ActionButtonColor.RED}
                    />
                </div>
            </div>
        </OverlayComponent>
    );
};

export default GeneratedInfoPopup;
