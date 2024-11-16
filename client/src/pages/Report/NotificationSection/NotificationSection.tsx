import SectionComponent from 'components/SectionComponent/SectionComponent';
import SVGIcon from 'components/SVGIcon/SVGIcon';
import React, { useState } from 'react';
import NotificationChannelTable from './NotificationChannelTable';
import OverlayComponent from 'components/OverlayComponent/OverlayComponent.tsx';
import NotificationsEntriesSelector
    from 'components/EntriesSelector/NotificationsEntriesSelector/NotificationsEntriesSelector.tsx';

export interface NotificationChannel {
    id: string;
    name: string;
    service: 'SLACK' | 'DISCORD' | 'EMAIL';
    details: string;
    updated: string;
    added: string;

    [key: string]: string;
}

interface NotificationSectionProps {
    notificationChannels: NotificationChannel[];
    setNotificationChannels: React.Dispatch<React.SetStateAction<NotificationChannel[]>>
}

const NotificationSection: React.FC<NotificationSectionProps> = ({
                                                                     setNotificationChannels,
                                                                     notificationChannels,
                                                                 }) => {
    const [selectedChannels, setSelectedChannels] = useState<NotificationChannel[]>([]);
    const [showModal, setShowModal] = useState(false);

    const handleAddClick = () => {
        setShowModal(true);
    };

    const handleCloseModal = () => {
        setShowModal(false);
    };

    const handleDelete = (id: string, service: string) => {
        setNotificationChannels(notificationChannels.filter(
            (channel) => !(channel.id === id && channel.service === service)
        ));
    };

    const handleAddSelected = () => {
        setNotificationChannels(prevChannels => {
            const newChannels = selectedChannels.filter(
                (channel) => !prevChannels.some(
                    (row) => row.id === channel.id && row.service === channel.service
                )
            );
            return [...prevChannels, ...newChannels];
        });
        setShowModal(false);
        setSelectedChannels([]);
    };

    return (
        <SectionComponent
            icon={<SVGIcon iconName="notification-icon" />}
            title="Notification channels"
            callback={handleAddClick}
        >
            <OverlayComponent isDisplayed={showModal} onClose={handleCloseModal}>
                <NotificationsEntriesSelector
                    selectedChannels={selectedChannels}
                    setSelectedChannels={setSelectedChannels}
                    channelsToExclude={notificationChannels}
                    onAdd={handleAddSelected}
                    onClose={handleCloseModal}
                />
            </OverlayComponent>
            {notificationChannels.length === 0 ? (
                <p>No notification channels selected, please add new.</p>
            ) : (
                <NotificationChannelTable rows={notificationChannels} onDelete={handleDelete} />
            )}
        </SectionComponent>
    );
};

export default NotificationSection;