import { UrgencyLevel } from 'api/managment-service';

export interface GenericIncident {
  // source: string;
  category: string;
  title: string;
  // timestamp: number;
  urgency: UrgencyLevel;
  // [index: string]: string | number;
}
