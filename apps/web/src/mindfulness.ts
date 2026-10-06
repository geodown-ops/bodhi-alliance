// 科技正念減壓八週課程。內容整理自佛教雲端大學《正念減壓課程》2022／2023 年版講義，法源法師審定編輯。
// 圖表改寫成文字；科學與健康說法附出處（refs 對應 references 的 key）。

export interface Ref {
  text: string
  url?: string
}

export const references: Record<string, Ref> = {
  kz1982: {
    text: 'Kabat-Zinn, J. (1982). An outpatient program in behavioral medicine for chronic pain patients based on the practice of mindfulness meditation. General Hospital Psychiatry, 4(1), 33–47.',
    url: 'https://doi.org/10.1016/0163-8343(82)90026-3',
  },
  kz1990: {
    text: 'Kabat-Zinn, J. (1990／2013 修訂版). Full Catastrophe Living. New York: Delacorte／Bantam.（正念飲食、身體掃描與七個態度；修訂版加入感恩與慷慨）',
  },
  kz2003: {
    text: 'Kabat-Zinn, J. (2003). Mindfulness-based interventions in context: past, present, and future. Clinical Psychology: Science and Practice, 10(2), 144–156.',
    url: 'https://doi.org/10.1093/clipsy.bpg016',
  },
  goyal2014: {
    text: 'Goyal, M., et al. (2014). Meditation programs for psychological stress and well-being: a systematic review and meta-analysis. JAMA Internal Medicine, 174(3), 357–368.',
    url: 'https://doi.org/10.1001/jamainternmed.2013.13018',
  },
  khoury2015: {
    text: 'Khoury, B., Sharma, M., Rush, S. E., & Fournier, C. (2015). Mindfulness-based stress reduction for healthy individuals: a meta-analysis. Journal of Psychosomatic Research, 78(6), 519–528.',
    url: 'https://doi.org/10.1016/j.jpsychores.2015.03.009',
  },
  jacobson1938: {
    text: 'Jacobson, E. (1938). Progressive Relaxation. Chicago: University of Chicago Press.（先緊後鬆的放鬆法）',
  },
  taskforce1996: {
    text: 'Task Force of the European Society of Cardiology and the North American Society of Pacing and Electrophysiology (1996). Heart rate variability: standards of measurement, physiological interpretation and clinical use. Circulation, 93(5), 1043–1065.',
    url: 'https://doi.org/10.1161/01.CIR.93.5.1043',
  },
  shaffer2017: {
    text: 'Shaffer, F., & Ginsberg, J. P. (2017). An overview of heart rate variability metrics and norms. Frontiers in Public Health, 5, 258.',
    url: 'https://doi.org/10.3389/fpubh.2017.00258',
  },
  hirsch1981: {
    text: 'Hirsch, J. A., & Bishop, B. (1981). Respiratory sinus arrhythmia in humans: how breathing pattern modulates heart rate. American Journal of Physiology, 241(4), H620–H629.',
    url: 'https://doi.org/10.1152/ajpheart.1981.241.4.H620',
  },
  lehrer2014: {
    text: 'Lehrer, P. M., & Gevirtz, R. (2014). Heart rate variability biofeedback: how and why does it work? Frontiers in Psychology, 5, 756.',
    url: 'https://doi.org/10.3389/fpsyg.2014.00756',
  },
  zaccaro2018: {
    text: 'Zaccaro, A., et al. (2018). How breath-control can change your life: a systematic review on psycho-physiological correlates of slow breathing. Frontiers in Human Neuroscience, 12, 353.',
    url: 'https://doi.org/10.3389/fnhum.2018.00353',
  },
  lomas2015: {
    text: 'Lomas, T., Ivtzan, I., & Fu, C. H. Y. (2015). A systematic review of the neurophysiology of mindfulness on EEG oscillations. Neuroscience & Biobehavioral Reviews, 57, 401–410.',
    url: 'https://doi.org/10.1016/j.neubiorev.2015.09.018',
  },
  russell1980: {
    text: 'Russell, J. A. (1980). A circumplex model of affect. Journal of Personality and Social Psychology, 39(6), 1161–1178.（情緒座標：激動程度 × 愉悅程度）',
    url: 'https://doi.org/10.1037/h0077714',
  },
  gross1998: {
    text: 'Gross, J. J. (1998). The emerging field of emotion regulation: an integrative review. Review of General Psychology, 2(3), 271–299.',
    url: 'https://doi.org/10.1037/1089-2680.2.3.271',
  },
  myers1985: {
    text: 'Myers, I. B., & McCaulley, M. H. (1985). Manual: A Guide to the Development and Use of the Myers-Briggs Type Indicator. Palo Alto: Consulting Psychologists Press.（依據 Jung, C. G. 1921《心理類型》）',
  },
  pittenger2005: {
    text: 'Pittenger, D. J. (2005). Cautionary comments regarding the Myers-Briggs Type Indicator. Consulting Psychology Journal: Practice and Research, 57(3), 210–221.',
    url: 'https://doi.org/10.1037/1065-9293.57.3.210',
  },
  goleman1995: {
    text: 'Goleman, D. (1995). Emotional Intelligence. New York: Bantam.（情商五能力）',
  },
  mcewen1998: {
    text: 'McEwen, B. S. (1998). Protective and damaging effects of stress mediators. New England Journal of Medicine, 338(3), 171–179.',
    url: 'https://doi.org/10.1056/NEJM199801153380307',
  },
  cohen2007: {
    text: 'Cohen, S., Janicki-Deverts, D., & Miller, G. E. (2007). Psychological stress and disease. JAMA, 298(14), 1685–1687.',
    url: 'https://doi.org/10.1001/jama.298.14.1685',
  },
  wright2002: {
    text: 'Wright, R. J., Cohen, S., Carey, V., Weiss, S. T., & Gold, D. R. (2002). Parental stress as a predictor of wheezing in infancy: a prospective birth-cohort study. American Journal of Respiratory and Critical Care Medicine, 165(3), 358–365.',
    url: 'https://doi.org/10.1164/ajrccm.165.3.2102016',
  },
  steptoe2012: {
    text: 'Steptoe, A., & Kivimäki, M. (2012). Stress and cardiovascular disease. Nature Reviews Cardiology, 9(6), 360–370.',
    url: 'https://doi.org/10.1038/nrcardio.2012.45',
  },
  mayer2000: {
    text: 'Mayer, E. A. (2000). The neurobiology of stress and gastrointestinal disease. Gut, 47(6), 861–869.',
    url: 'https://doi.org/10.1136/gut.47.6.861',
  },
  epel2000: {
    text: 'Epel, E. S., et al. (2000). Stress and body shape: stress-induced cortisol secretion is consistently greater among women with central fat. Psychosomatic Medicine, 62(5), 623–632.',
    url: 'https://doi.org/10.1097/00006842-200009000-00005',
  },
  surwit1992: {
    text: 'Surwit, R. S., Schneider, M. S., & Feinglos, M. N. (1992). Stress and diabetes mellitus. Diabetes Care, 15(10), 1413–1422.',
    url: 'https://doi.org/10.2337/diacare.15.10.1413',
  },
  wilson2003: {
    text: 'Wilson, R. S., et al. (2003). Proneness to psychological distress is associated with risk of Alzheimer’s disease. Neurology, 61(11), 1479–1485.',
    url: 'https://doi.org/10.1212/01.WNL.0000096167.56734.59',
  },
  epel2004: {
    text: 'Epel, E. S., et al. (2004). Accelerated telomere shortening in response to life stress. PNAS, 101(49), 17312–17315.（壓力最高者的端粒，相當於多老化約 9–17 年）',
    url: 'https://doi.org/10.1073/pnas.0407162101',
  },
  hawkins1995: {
    text: 'Hawkins, D. R. (1995). Power vs. Force: The Hidden Determinants of Human Behavior. Sedona: Veritas.（中譯《心靈能量》）',
  },
  satipatthana: {
    text: '《中阿含經》卷 24〈念處經〉（T01, no. 26）；巴利《念處經》（Satipaṭṭhāna Sutta, MN 10）。四念處：身、受、心、法。',
  },
  anapana: {
    text: '《雜阿含經》卷 29（T02, no. 99，第 810 經等）；巴利《入出息念經》（Ānāpānasati Sutta, MN 118）。安般念十六勝行。',
  },
  visuddhimagga: {
    text: '覺音《清淨道論》（Visuddhimagga）第八品〈隨念業處〉入出息念。',
  },
  zhiyi: {
    text: '智顗《釋禪波羅蜜次第法門》（T46, no. 1916）、《六妙法門》（T46, no. 1917）。數、隨、止、觀、還、淨。',
  },
  sallatha: {
    text: '《雜阿含經》卷 17 第 470 經（T02, no. 99）；巴利《箭經》（Sallatha Sutta, SN 36.6）。三受與「第二支箭」。',
  },
  dajifa: {
    text: '《大集法門經》卷 1（T01, no. 12, p. 228b29–c7）。四禪。',
  },
}

export const attitudes = [
  { name: '不作判斷', en: 'Non-Judging', text: '不替自己的情緒、想法、病痛等身心現象打分數，只是純粹地覺察它們。' },
  { name: '保持耐心', en: 'Patience', text: '對自己當下的各種身心狀況保持耐心，和它們和平共處。' },
  { name: '初學之心', en: "Beginner's Mind", text: '常保初學者的心，以赤子之心面對每一個身心事件。' },
  { name: '信任自己', en: 'Trust', text: '相信自己的智慧與能力。' },
  { name: '不要強求', en: 'Non-Striving', text: '只是覺察當下發生的一切，不強求想要的結果。' },
  { name: '接受現狀', en: 'Acceptance', text: '願意如實地觀照自己此刻的身心。' },
  { name: '放下種種', en: 'Letting Go', text: '放下好惡分別，分分秒秒覺察當下的身心事件。' },
  { name: '感恩一切', en: 'Gratitude', text: '感謝此刻擁有的生命與一切，允許自己好好接受它。' },
  { name: '慷慨布施', en: 'Generosity', text: '把自己的時間、關懷與所學，慷慨地分享出去。' },
]

export const dedication = [
  '願我無身體上的痛苦',
  '願我無精神上的痛苦',
  '願我安祥與快樂',
  '願一切大眾無身體上的痛苦',
  '願一切大眾無精神上的痛苦',
  '願大家安祥與快樂',
]

export interface Section {
  title: string
  paras?: string[]
  items?: string[]
  note?: string
  refs?: string[]
}

export interface Week {
  n: number
  part: '身念住' | '受念住' | '心念住' | '法念住'
  title: string
  summary: string
  sections: Section[]
  steps: string[]
  life: string
}

export const parts = [
  { name: '身念住', text: '從身體與呼吸開始，讓心安定下來。', weeks: [1, 2] },
  { name: '受念住', text: '看清苦、樂、不苦不樂三種感受。', weeks: [3, 4] },
  { name: '心念住', text: '認識情緒與性格，學會調節。', weeks: [5, 6] },
  { name: '法念住', text: '覺知無常，與自己、他人、人生和解。', weeks: [7, 8] },
]

export const weeks: Week[] = [
  {
    n: 1,
    part: '身念住',
    title: '認識正念',
    summary: '正念是什麼、要用什麼態度練，從喝一口水、吃一顆葡萄乾和身體掃描開始。',
    sections: [
      {
        title: '正念減壓從哪裡來',
        paras: [
          '喬．卡巴金博士（Jon Kabat-Zinn，1944 年生，麻省理工學院分子生物學博士）於 1979 年在麻州大學醫學院創立正念減壓（MBSR）課程與正念中心，是把佛教正念禪修帶進西方主流醫學的第一人。',
          '他對正念的說法是：刻意專注於當下的一種覺察，暫時收起判斷心，帶著好奇與接納，一刻接著一刻。',
        ],
        refs: ['kz1982', 'kz2003'],
      },
      {
        title: '何謂正念',
        paras: [
          '「念」拆開是「今＋心」：此刻心中的念頭、想法、感受與情緒。「正」是當下、不偏不倚、沒有對立，不是「正向」或「正確」的意思。',
          '所以不論正向或負向的念頭，都是被觀察、認識與接納的對象。正念就是時時刻刻活在當下的能力，讓我們活得了了分明、自主自在。',
        ],
      },
      {
        title: '課程目標：離苦得樂',
        paras: ['正念減壓一路從放鬆抒壓，到輕安喜悅，再到安穩快樂；在佛法上，是從正念、正定、正知走向覺悟解脫。'],
      },
      {
        title: '科技正念：看見練習的效果',
        paras: [
          '課程會用心律變異（HRV）與腦波量測，讓學員看見練習前後的變化。心律量測看平均心律、心律變異度，以及代表副交感神經（放鬆）與交感神經（緊張）活性的頻譜指標；腦波帶透過手機 App，把當下的專注力與放鬆力換算成百分比。',
          '這些數值只用來觀察自己的練習，不是醫療檢查。',
        ],
        refs: ['taskforce1996', 'shaffer2017', 'lomas2015'],
      },
      {
        title: '身體掃描：先緊後鬆',
        paras: ['從頭到腳一處一處覺察，每一處先用力繃緊幾秒，再完全放鬆，感受兩者的差別：眼睛、臉頰、後頸，聳肩、壓背、挺腰，最後提腳跟放下雙腿。'],
        refs: ['jacobson1938', 'kz1990'],
      },
    ],
    steps: [
      '正念喝水：準備一杯水，先用手感覺杯子的溫度與重量，再看水的樣子、聞一聞，最後小口慢慢喝，感覺水流過口腔與喉嚨。',
      '正念飲食：準備五顆葡萄乾或其他果乾，一顆一顆依序用觸覺、視覺、嗅覺、味覺去認識它，慢慢咀嚼，注意滋味的變化。',
      '身體掃描：坐好或躺好，從眼睛、臉頰、後頸、肩膀、背、腰到雙腿，每一處先繃緊再放鬆。',
    ],
    life: '每天找時間做身體掃描；吃飯、喝茶、穿衣、洗澡時，用正念飲食的方式覺察。',
  },
  {
    n: 2,
    part: '身念住',
    title: '正念呼吸',
    summary: '腹式呼吸與數息、隨息。慢慢吐氣，讓身體自己放鬆下來。',
    sections: [
      {
        title: '安般念',
        paras: ['安般念（Ānāpāna，安那般那念）就是念出入息：清楚地觀察、覺知自己的呼吸。這週練腹式呼吸，以及數息和隨息。'],
        refs: ['anapana'],
      },
      {
        title: '腹式呼吸',
        paras: [
          '胸式呼吸短而急；腹式呼吸既長且深，吸氣時腹部像氣球一樣鼓起，吐氣時收縮。它能吸進更多氧氣，讓人感到放鬆、安定神經、改善專注，並幫助自律神經平衡。',
        ],
      },
      {
        title: '呼吸可以調節心律',
        paras: [
          '平常每分鐘呼吸 12 到 18 次時，心跳不太跟著呼吸變；呼吸放慢後，心跳開始跟著呼吸起伏：吸氣時變快，吐氣時變慢，這叫呼吸性竇性心律不整（RSA）。',
          '吸氣時胸腔壓力降低、吐氣時升高，透過壓力感受器反射與迷走神經傳到腦幹的孤束核，再影響掌管情緒的杏仁核與下視丘。所以：放慢吸氣偏向提振交感神經（提神、應戰），放慢吐氣偏向提升副交感神經（放鬆、好眠）。',
        ],
        refs: ['hirsch1981', 'lehrer2014', 'zaccaro2018'],
      },
      {
        title: '數息觀與六妙門',
        paras: [
          '完整的數息觀有六門：數、隨、止、觀、還、淨。數息、隨息、止三門，從正念練到正定；觀、還、淨三門，從正定練到正知，這也是世間禪與出世間禪的分界。數息是入門的第一步。',
          '對照南傳《清淨道論》、天台《釋禪波羅蜜》與淨土念佛三昧，這一段屬於身念住：覺知息長、息短、全身息，到安息身行，效果是紓壓，對治粗心與昏沉散亂。',
        ],
        refs: ['visuddhimagga', 'zhiyi'],
      },
    ],
    steps: [
      '先做幾次腹式呼吸：鼻子吸氣、腹部鼓起；鼻子吐氣、腹部收縮，慢慢放慢。',
      '數息：每完成一次完整的吸氣和吐氣，心裡數一；再一次數二，數到十再從一開始，十個一組，持續循環。',
      '數到心初步靜下來，就轉成隨息：不再數，只是清楚覺知整個呼吸的過程，從吸進到吐完。',
      '如果昏沉或散亂，就回到數息，直到心不動亂、身體安穩。',
    ],
    life: '每天找時間做正念呼吸；把呼吸時的覺察帶進吃飯、喝茶、穿衣、洗澡。',
  },
  {
    n: 3,
    part: '受念住',
    title: '三種感受',
    summary: '苦受、樂受、不苦不樂受。看見感受，就不會被它牽著走。',
    sections: [
      {
        title: '三種覺受',
        items: [
          '苦受：痛苦的感受。失去覺知時容易引起瞋恚，帶來煩惱。',
          '樂受：快樂的感受。失去覺知時容易引起貪欲，滿足後又煩躁、不滿足。',
          '不苦不樂受：例如無聊。失去覺知時容易引起愚癡，想找東西來餵養自己的心，心就靜不下來。',
        ],
        paras: ['有正念覺知時，樂受帶來滿足、苦受帶來平靜，不苦不樂受也只是平靜的一刻。'],
        refs: ['sallatha'],
      },
      {
        title: '覺受原則',
        items: [
          '就感受觀察感受，不在感受之外繞圈子。',
          '從內部、外部、內外部三方面觀察感受。',
          '觀察感受的生起、滅去，體會無常變化。',
          '感受只是感受，其中沒有「我」，也不打分數。',
        ],
      },
      {
        title: '在修行地圖上的位置',
        paras: ['這是受念住：覺知喜、覺知樂、覺知心行、覺知安息心，止觀上從念而無念進入初禪、二禪。'],
        refs: ['satipatthana', 'anapana'],
      },
    ],
    steps: [
      '先數息，再隨息，讓心靜下來。',
      '身體受念：用身體掃描的方法，觀察身上哪裡是苦受、哪裡是樂受、哪裡是不苦不樂受。',
      '心理受念：覺知全身的呼吸，練到不昏沉、不散亂，再觀察心裡的苦受、樂受、不苦不樂受。',
    ],
    life: '每天找時間做受念覺知；把對感受的覺察帶進吃飯、喝茶、穿衣、洗澡。',
  },
  {
    n: 4,
    part: '受念住',
    title: '靜坐與經行',
    summary: '調身、調息、調心：學會坐得穩，也學會帶著覺察走路。',
    sections: [
      {
        title: '靜坐：調身、調息、調心',
        paras: [
          '調身用七支坐法；調息用腹式呼吸，又叫丹田呼吸，想像把空氣一路吸到下腹；調心是訓練專注，同時也放鬆。',
        ],
        items: [
          '腿：雙盤、單盤、散盤或坐椅子都可以，重點是坐穩。',
          '手：結禪定印，雙手上下交疊，兩個拇指輕輕相觸。',
          '身：脊背正直，不彎腰駝背。',
          '肩：雙肩平衡放鬆。',
          '下巴：微收，後頸拉直。',
          '舌：舌尖輕抵上顎。',
          '眼：開三分、閉七分。',
        ],
      },
      {
        title: '經行（行禪）',
        paras: [
          '經行是在一段直線上來回慢慢走，帶著覺性走路。全身放鬆，雙手輕抱胸前或握在身後，不東張西望，也不默數默念，只是覺知腳的移動：提起、往前、放下、著地，越慢越清楚。',
          '經行是動態的禪定訓練。行住坐臥、動靜語默，做任何事都可以練習帶著覺性。',
        ],
      },
    ],
    steps: [
      '用七支坐法坐好，做幾次腹式呼吸。',
      '先數息、再隨息，延續上週的身體受念與心理受念。',
      '起身經行十分鐘：一步一步，覺知提起、往前、放下、著地。',
      '完整的上座與下座方法，請看「上座與下座」。',
    ],
    life: '每天找時間靜坐或經行；把受念的覺察帶進吃飯、喝茶、穿衣、洗澡。',
  },
  {
    n: 5,
    part: '心念住',
    title: '情緒與心念',
    summary: '情緒沒有好壞，它只是在傳遞訊息。先找到它在情緒座標上的位置，再用呼吸調節。',
    sections: [
      {
        title: '情緒是什麼',
        paras: [
          '情緒是感覺、思維和行為綜合起來的心理產物，來自外境與內心的交互作用，讓我們趨吉避凶。不論消極或積極，情緒的出現都正常，只是在告訴我們：此刻的情境讓人不愉快、焦慮、害怕，或愉快、興奮。',
        ],
        refs: ['gross1998'],
      },
      {
        title: '情緒座標',
        paras: ['用兩條軸來定位情緒：上下是激動或不激動，左右是愉悅或不愉悅。'],
        items: [
          '激動又不愉悅：驚慌、發怒、懊惱、挫折。',
          '激動又愉悅：興奮、驚喜、高興、快樂。',
          '不激動也不愉悅：悲傷、沮喪、無聊、疲倦。',
          '不激動而愉悅：安詳、平靜、放鬆，再往下就是昏沉。',
        ],
        refs: ['russell1980'],
      },
      {
        title: '情緒的本質',
        items: [
          '情緒沒有好壞之分，但需要適度調節，以免失控。',
          '真正的困擾不是負面情緒本身，而是對待情緒的態度和行為。',
          '大多數人的困擾，來自不能接納當下的情緒，越想「解決情緒」越陷進去。',
          '缺少調節的方法，在強烈情緒裡就容易失控。',
          '用正念的方式調節情緒，就能消除情緒帶來的不愉快，活得更積極。',
        ],
      },
      {
        title: '在修行地圖上的位置',
        paras: ['從受念住進入心念住：覺知心、令心喜悅、令心等持、令心解脫，止觀上進入觀門，對治細微的心念。'],
        refs: ['satipatthana'],
      },
    ],
    steps: [
      '先做正念呼吸，讓心靜下來。',
      '覺知此刻的情緒和心念，在情緒座標上找到它的位置。',
      '觀想呼吸：吸氣時想像能量進入身體，吐氣時想像情緒排出體外。',
    ],
    life: '每天找時間覺知情緒與心念；把心念的覺察帶進吃飯、喝茶、穿衣、洗澡。',
  },
  {
    n: 6,
    part: '心念住',
    title: '性格與心理',
    summary: '了解自己和別人的性格，接納情緒，不讓它累積成情結。',
    sections: [
      {
        title: '性格的四種特徵',
        paras: ['性格表現在我們對自己、對別人、對事物的態度與言行上，是個性的核心。'],
        items: [
          '態度特徵：誠實或虛偽、謙遜或驕傲。',
          '意志特徵：勇敢或怯懦、果斷或優柔寡斷。',
          '理智特徵：敏捷或遲緩、深刻或淺薄、有無邏輯。',
          '情緒特徵：熱情或冷漠、開朗或抑鬱。',
        ],
      },
      {
        title: 'MBTI 性格指標',
        paras: [
          'MBTI 由布里格斯與邁爾斯母女依榮格的心理類型理論發展而成，用四個維度描述性格：',
        ],
        items: [
          '外向 E／內向 I：從他人身上得到動力、喜歡快步調；或喜歡獨自慢慢來、一次專注一件事。',
          '實感 S／直覺 N：重事實細節與過去經驗；或看見可能性與遠景、喜歡創新。',
          '理智 T／感性 F：依客觀事實與分析做決定；或從個人價值出發、關心他人感受。',
          '決斷 J／彈性 P：喜歡有條理、按計畫；或不介意突發狀況、喜歡彈性。',
        ],
        note: 'MBTI 適合用來認識自己與別人的差異；學術界對它的信度與效度仍有討論，不宜用來給人貼標籤。',
        refs: ['myers1985', 'pittenger2005'],
      },
      {
        title: '避免形成情結',
        paras: ['情緒一再累積，會變成情結，再形成惡性循環。透過了解自己與他人的性格，全然接納情緒，就能從情結裡解脫出來。'],
      },
      {
        title: '培養情商的五種能力',
        items: [
          '自我覺察：了解自己感受的能力。',
          '自我管理：處理自己情緒與衝動的能力。',
          '自我激勵：面對失敗挫折仍屹立不搖的能力。',
          '同理心：體會與了解他人感受的能力。',
          '社交能力：善於對待與處理他人情緒的能力。',
        ],
        refs: ['goleman1995'],
      },
    ],
    steps: [
      '先透過 MBTI 認識自己的性格傾向。',
      '做正念呼吸，讓心靜下來。',
      '覺知此刻的情緒與心念，在情緒座標上定位。',
      '觀想呼吸：吸氣時想像能量進入，吐氣時想像情緒排出。',
    ],
    life: '每天找時間覺知情緒與心念；把心念的覺察帶進吃飯、喝茶、穿衣、洗澡。',
  },
  {
    n: 7,
    part: '法念住',
    title: '專注與放鬆',
    summary: '既專注又放鬆，自律神經才平衡。覺知生滅無常，舒緩外在與內心的壓力。',
    sections: [
      {
        title: '放鬆式專注',
        paras: [
          '只有專注、只提高交感神經，身心長期緊繃，壓力就會累積，就像專心上班一整天，下班後特別累。只有放鬆、只提高副交感神經，長久下來反而懈怠，就像放假在家休息太久，反而倦怠。',
          '正念練習帶來的專注不一樣：在專注的同時也覺察到身體的緊繃，並讓它放鬆。持續練習，就能同時專注又放鬆，自律神經平衡，紓解壓力又身心平衡。把這份覺察帶進生活、工作、家庭與人際，就能抒解壓力、發揮潛能，活得更輕鬆。',
        ],
        refs: ['lehrer2014', 'mcewen1998'],
      },
      {
        title: '腦波',
        paras: [
          '腦波儀（EEG）量測大腦的電位變化，依頻率分為 β 波（約 12–38 Hz，清醒專注）、α 波（約 8–12 Hz，放鬆清醒）、θ 波（約 4–8 Hz，深度放鬆、冥想）與 δ 波（約 0.5–4 Hz，深睡）。研究整理發現，正念練習時 α 波與 θ 波的能量常會增加。',
        ],
        refs: ['lomas2015'],
      },
      {
        title: '內外壓力',
        items: [
          '外來壓力：家庭狀況、身邊的人、工作量與工作變化、經濟壓力等。',
          '內在壓力：要求完美、排不出優先順序、追逐目標、內心對老死的恐懼等。',
        ],
      },
      {
        title: '壓力與疾病',
        paras: [
          '依衛福部資料，與壓力有關的疾病遍及神經（偏頭痛、緊張性頭痛、背痛）、內分泌（月經不規律）、消化（潰瘍、腸道發炎）、呼吸（氣喘、花粉熱）、心血管（高血壓、中風、冠心病）、生殖與免疫系統（濕疹、蕁麻疹、乾癬、過敏），睡眠問題、肩頸痠痛、皮膚問題也常和壓力有關。',
        ],
        items: [
          '氣喘：父母長期的壓力，與孩子日後喘鳴、氣喘的風險有關。',
          '心血管：突發的強烈情緒壓力可能誘發心血管事件；長期壓力也是心臟病的風險因子。',
          '腸胃：壓力和胃食道逆流、大腸激躁症密切相關。',
          '肥胖：壓力荷爾蒙與腹部肥胖有關。',
          '糖尿病：壓力會升高血糖，第二型糖尿病患者尤其明顯。',
          '阿茲海默症：容易感到心理壓力的人，罹病風險較高。',
          '老化：長期高壓的人，細胞的端粒明顯較短，相當於多老化約九到十七年。',
        ],
        note: '這些是研究中觀察到的關聯，不代表壓力是唯一原因；身體不適請先就醫。',
        refs: ['cohen2007', 'mcewen1998', 'wright2002', 'steptoe2012', 'mayer2000', 'epel2000', 'surwit1992', 'wilson2003', 'epel2004'],
      },
      {
        title: '覺知無常',
        paras: [
          '隨息時，心隨著呼吸出入，清楚知道每一口氣的遠近長短。觀察呼吸一生一滅，體會苦、無常、無我，能舒緩外界情境造成的壓力；再體會那不生不滅、真常的覺性，能舒緩內心對老死的恐懼。',
          '這是法念住：觀無常、觀離欲、觀滅、觀捨遣，止觀雙運、定慧等持，增上智慧。',
        ],
        refs: ['anapana', 'zhiyi'],
      },
    ],
    steps: [
      '做正念呼吸，讓心靜下來。',
      '覺知每一口呼吸的生與滅，體會無常，放下外在情境帶來的壓力。',
      '再覺知那不生不滅的覺性，放下對老死的恐懼。',
      '觀想呼吸：吸氣時想像能量進入身體，吐氣時想像情緒排出體外。',
    ],
    life: '每天找時間練放鬆式專注；把無常的覺察帶進吃飯、喝茶、穿衣、洗澡。',
  },
  {
    n: 8,
    part: '法念住',
    title: '和解與快樂',
    summary: '與自己、與他人、與人生和解。活在當下、接納當下，快樂自然生起。',
    sections: [
      {
        title: '三種和解',
        items: [
          '與自己和解：正念讓我們看見自己學來的自我否定，看見了才能改掉，才能與自己和解。',
          '與他人和解：別人的傷害，多半來自他們的無意識、恐懼、制約反應與經歷；我們也會因無意識而傷害別人。心懷怨恨，受苦的是自己。',
          '與人生和解：我們常苦於希望人生不是這個樣子。與擁有的一切和解，才會真正感恩自己的人生。',
        ],
      },
      {
        title: '快樂的要件',
        items: [
          '與自己連結：從覺察身體當下的感覺開始，呼吸、吃東西、走路、伸展都是練習；和身體連結的人，更懂得珍惜與穩住自己。',
          '與他人連結：練習不打岔地專注聆聽、向對方核對自己的理解，有覺察地說出意見，而不是出於慣性。',
          '照顧好自己：覺察自己真實的狀態與需求，適度地滿足它。',
          '看到整體脈絡：答案不明時人容易焦躁，越動越亂；先穩住自己，才看得到全局。',
          '在覺察中行動：覺察行動是否真正有益、時機是否恰當；不因害怕或拖延而不動，也不在厭惡、忌妒中行動。',
        ],
      },
      {
        title: '快樂的層次',
        paras: [
          '正念練習友善地對待自己、別人與眾生。覺察帶來看見，看見帶來選擇，明智的選擇帶來快樂：從放鬆喜樂、輕安悅樂，到平等喜捨、寂靜妙樂；在修行上，是正念、正定、正智到正解脫。',
        ],
      },
      {
        title: '霍金斯的意識能級',
        paras: [
          '美國的霍金斯（David R. Hawkins）把人類的意識狀態對應到 1 到 1000 的能量刻度，分成十七個能級。高能量由上而下是：開悟正覺（700–1000）、寧靜極樂（600）、平和喜悅（540）、仁愛崇敬（500）、理解諒解（400）、寬容接納（350）、主動樂觀（310）、信任淡定（250）、勇氣肯定（200）；低能量是：驕傲刻薄（175）、憤怒仇恨（150）、欲望渴求（125）、恐懼焦慮（100）、憂傷無助（75）、冷漠絕望（50）、內疚報復（30）、羞恥蔑視（20 以下）。',
          '課程的重點不在刻意追求高能級：清楚地活在當下、接納當下所有感受與情緒，超越二元對立的評判，內在的正能量就自然生起，人也輕鬆自在。遇到劇烈的情緒或重大煩惱，則可以用觀想來得到平復與療癒。',
        ],
        note: '霍金斯的能級是以肌肉測試的方法得出，屬於他個人的研究與著作，尚未經過主流科學驗證。',
        refs: ['hawkins1995'],
      },
    ],
    steps: [
      '做正念呼吸，讓心靜下來。',
      '覺知生滅無常，體會煩惱終將消除；覺知不生滅的真常，體會快樂終將到來。',
      '觀想呼吸：吸氣時想像聖賢、覺者或宇宙的光明能量進入身體，喚醒內在的能量；吐氣時想像心中充滿快樂，身心清淨平衡。',
    ],
    life: '每天找時間做正念增能練習；把對快樂的覺察帶進吃飯、喝茶、穿衣、洗澡。',
  },
]

// 上座與下座，依法源法師〈止觀禪修上下坐引導〉逐字稿整理
export const sitting = {
  before: ['上座前先調飲食、調睡眠：不要太飽或太餓，也不要睡眠不足。'],
  up: [
    {
      title: '調身：七支坐法',
      items: [
        '腿：可以盤腿（雙盤、單盤、散盤），也可以坐椅子，重點是坐穩。',
        '手：結禪定印，雙手上下交疊，拇指輕輕相觸。',
        '身：保持正直。先把上身往前傾，讓坐骨往後打開，再從頭、頸、背、腰一節一節慢慢坐起來；坐骨像三角架撐住身體，不必提肛或腰部用力。',
        '肩：吸氣聳肩，吐氣把力量放掉，做兩次，肩膀就放鬆了。',
        '下巴：微收，後頸拉直，頭部氣血通暢。',
        '舌：舌頂上顎，不會口乾舌燥。',
        '眼：開三分、閉七分。先看向前下方 45 度，眼皮垂下只留三成縫隙，再把視線抬平；原本看到的部分變成餘光，就不會干擾。',
      ],
    },
    {
      title: '調息：腹式呼吸',
      items: ['鼻吸鼻吐。吸氣時腹部像吹氣球一樣鼓起，吐氣時收縮；呼吸放慢，想像吸進清新的空氣，吐出體內的濁氣。調好呼吸，心就比較容易靜下來。'],
    },
    {
      title: '調心：觀照',
      items: [
        '選一個專注的對象：呼吸，或四字、六字佛號。',
        '先練到不散亂、不昏沉，再清楚觀察它的生滅變化，觀照無常。慢慢地，身體會記得「無我」的感覺，遇到任何情境就比較不煩惱、不執著。',
        '過程中有昏沉或散亂就盡量對治；真的很睏，就告訴自己睡三分鐘、五分鐘，醒來繼續。',
      ],
    },
  ],
  down: [
    { title: '調心', items: ['先在心裡準備下座；如果有入定，就出定、捨掉禪相。'] },
    { title: '調息', items: ['呼吸不順就多坐一下，做幾次深長的腹式呼吸。不要馬上站起來。'] },
    {
      title: '調身：收功動作',
      items: [
        '身體前後左右輕輕搖晃，從靜慢慢轉為動。',
        '熨目：雙手搓熱，掌心輕蓋眼球，讓眼睛吸收手的溫度，再拿開、睜眼。',
        '搓臉：雙手搓熱，中指沿鼻樑旁、眼眶向外畫圓，把臉上的氣搓勻。',
        '梳頭：十指從前額往後腦梳，按摩頭頂。',
        '搓後頸，再雙手放在肩井穴，由後往前拉一拉肩膀。',
        '熨腎、擦腎：雙手搓熱貼在後腰；再用中指找脊柱兩側的凹陷，手掌貼著往下推。',
        '拍手臂：一手從上往下拍另一手的外側與內側，再換手。',
        '按摩膝蓋、腳踝等痠麻的地方；雙腳往斜前方伸展，用腰腹的力量做長短腿，舒緩腰與髖關節。',
        '雙手握空拳，敲大腿到小腿的外側（膽經）與內側（肝經），下座後多喝點溫開水。',
      ],
    },
  ],
  table: [
    { k: '調身', up: '七支坐法', mid: '保持正直', down: '收功動作' },
    { k: '調息', up: '腹式呼吸（吐濁納清）', mid: '自然呼吸', down: '自然呼吸' },
    { k: '調心', up: '由動到靜（數、隨、止）', mid: '置心一處（觀、還、淨）', down: '由靜到動' },
  ],
}
