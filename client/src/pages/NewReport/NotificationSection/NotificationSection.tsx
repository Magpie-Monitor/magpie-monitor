// NotificationSection.tsx
import './NotificationSection.scss';
import SectionComponent from 'components/SectionComponent/SectionComponent';
import SVGIcon from 'components/SVGIcon/SVGIcon';
import { useState, useEffect } from 'react';
import NotificationChannelTable from './NotificationChannelTable';
import OverlayComponent from 'components/OverlayComponent/OverlayComponent.tsx';
import { NotificationChannel } from 'api/managment-service';

const MOCK_CHANNELS: NotificationChannel[] = [
    { id: '1', name: 'Infra team slack', service: 'SLACK', details: 'wms_dev/#infra-alerts', updated: '07.03.2024 15:32', added: '07.03.2024 15:32' },
    { id: '2', name: 'Infra team discord', service: 'DISCORD', details: 'wms_dev/#dev-infra-alerts', updated: '07.03.2024 15:32', added: '07.03.2024 15:32' },
    { id: '3', name: 'Kontakt wms', service: 'EMAIL', details: 'kontakt@wmsdev.pl', updated: '07.03.2024 15:32', added: '07.03.2024 21:37' },
];

const NotificationSection = () => {
    const [rows, setRows] = useState<NotificationChannel[]>([]);
    const [loading, setLoading] = useState(true);
    const [showModal, setShowModal] = useState(false);

    useEffect(() => {
        setLoading(true);
        setTimeout(() => {
            setRows(MOCK_CHANNELS);
            setLoading(false);
        }, 2000);
    }, []);

    const handleAddClick = () => {
        setShowModal(true);
    };

    const handleCloseModal = () => {
        setShowModal(false);
    };

    return (
        <SectionComponent
            icon={<SVGIcon iconName="notification-icon" />}
            title="Notification Channels"
            callback={handleAddClick}
        >
            {showModal && <OverlayComponent onClose={handleCloseModal} />}
            <div className="notification-section__content">
                {loading ? (
                    <p>Loading...</p>
                ) : rows.length === 0 ? (
                    <p>No notification channels available.</p>
                ) : (
                    <NotificationChannelTable rows={rows} />
                )}
            </div>
        </SectionComponent>
    );
};

export default NotificationSection;
