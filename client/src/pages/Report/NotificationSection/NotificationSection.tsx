import SectionComponent from 'components/SectionComponent/SectionComponent';
import SVGIcon from 'components/SVGIcon/SVGIcon';
import React, { useState } from 'react';
import NotificationChannelTable from './NotificationChannelTable';
import OverlayComponent from 'components/OverlayComponent/OverlayComponent.tsx';
import NotificationsEntriesSelector
    from 'components/EntriesSelector/NotificationsEntriesSelector/NotificationsEntriesSelector.tsx';
import {ManagmentServiceApiInstance} from 'api/managment-service.ts';
import {dateFromTimestampMs} from 'lib/date.ts';

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
    const [availableChannels, setAvailableChannels] = useState<NotificationChannel[]>([]);

    const loadNotificationChannels = async () => {
        try {
            const channels =
              await ManagmentServiceApiInstance.getNotificationChannels();
            setAvailableChannels(channels.map(
              (channel): NotificationChannel => ({
                  id: channel.id,
                  name: channel.name,
                  service: channel.service,
                  details: channel.details,
                  updated: dateFromTimestampMs(channel.updated),
                  added: dateFromTimestampMs(channel.added),
              }),
            ))
        } catch (error) {
            console.error('Failed to fetch channels:', error);
        }
    };

    const handleOpenModal = async () => {
        await loadNotificationChannels();
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
            callback={handleOpenModal}
        >
            <OverlayComponent isDisplayed={showModal} onClose={handleCloseModal}>
                <NotificationsEntriesSelector
                    selectedChannels={selectedChannels}
                    setSelectedChannels={setSelectedChannels}
                    channelsToExclude={notificationChannels}
                    onAdd={handleAddSelected}
                    onClose={handleCloseModal}
                    availableChannels={availableChannels}
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