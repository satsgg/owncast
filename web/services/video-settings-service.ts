import { createContext } from 'react';

export type VideoQuality = {
  index: number;
  /**
   * This property is not just for display or so
   * but it holds information
   *
   * @example '1.2Mbps@24fps'
   */
  name: string;
  price: number;
  bandwidth: number;
};

export interface VideoSettingsStaticService {
  getVideoQualities(): Promise<Array<VideoQuality>>;
}

class VideoSettingsService {
  // TODO: Add price to response
  private static readonly VIDEO_CONFIG_URL = '/api/video/variants';

  public static async getVideoQualities(): Promise<Array<VideoQuality>> {
    let qualities: Array<VideoQuality> = [];

    try {
      const response = await fetch(VideoSettingsService.VIDEO_CONFIG_URL);
      qualities = await response.json();
    } catch (e) {
      console.error(e);
    }
    return qualities;
  }
}

export const VideoSettingsServiceContext =
  createContext<VideoSettingsStaticService>(VideoSettingsService);
