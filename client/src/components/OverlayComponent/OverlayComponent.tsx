import React from 'react';
import './OverlayComponent.scss';
import ActionButton, {ActionButtonColor} from 'components/ActionButton/ActionButton.tsx';

interface AddNotificationChannelModalProps {
    onClose: () => void;
}

const OverlayComponent: React.FC<AddNotificationChannelModalProps> = ({ onClose }) => {
    return (
        <div className="modal-overlay">
            <div className="modal-content">
                <p>List placeholder</p>
                <ActionButton
                    onClick={onClose}
                    description="Close"
                    color={ActionButtonColor.RED}
                />
            </div>
        </div>
    );
};

export default OverlayComponent;
