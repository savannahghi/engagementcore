package otp

// SendOtpToEmailTemplate generates an email template
const SendOtpToEmailTemplate = `

<!DOCTYPE html>
<html>

<head>
  <title>Be.Well by Slade 360°</title>
  <meta property="description" content="Be.Well One Time Pin" />
  <!--VIEWPORT-->
  <meta name="viewport" content="width=device-width; initial-scale=1.0; maximum-scale=1.0; user-scalable=no;" />
  <meta name="viewport" content="width=600, initial-scale = 2.3, user-scalable=no" />
  <meta name="viewport" content="width=device-width" />
  <!--CHARSET-->
  <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />

  <!-- IE=Edge and IE=X -->
  <meta http-equiv="X-UA-Compatible" content="IE=7" />
  <meta http-equiv="X-UA-Compatible" content="IE=8" />
  <meta http-equiv="X-UA-Compatible" content="IE=9" />
  <meta http-equiv="X-UA-Compatible" content="IE=edge" />
  <!--[if !mso]>-->
  <meta http-equiv="X-UA-Compatible" content="IE=edge" />
  <!--<![endif]-->
  <!-- INLINE STYLES -->
  <style type="text/css">
    @media screen {
      @font-face {
        font-family: "Lato";
        font-style: normal;
        font-weight: 400;
        src: local("Lato Regular"), local("Lato-Regular"),
          url(https://fonts.gstatic.com/s/lato/v11/qIIYRU-oROkIk8vfvxw6QvesZW2xOQ-xsNqO47m55DA.woff) format("woff");
      }

      @font-face {
        font-family: "Lato";
        font-style: normal;
        font-weight: 700;
        src: local("Lato Bold"), local("Lato-Bold"),
          url(https://fonts.gstatic.com/s/lato/v11/qdgUG4U09HnJwhYI-uK18wLUuEpTyoUstqEm5AMlJo4.woff) format("woff");
      }

      @font-face {
        font-family: "Lato";
        font-style: italic;
        font-weight: 700;
        src: local("Lato Bold Italic"), local("Lato-BoldItalic"),
          url(https://fonts.gstatic.com/s/lato/v11/HkF_qI1x_noxlxhrhMQYELO3LdcAZYWl9Si6vvxL-qU.woff) format("woff");
      }
    }

    body,
    table,
    td {
      -webkit-text-size-adjust: 100%;
      -ms-text-size-adjust: 100%;
      border: none !important;
    }

    table {
      border-collapse: collapse !important;
    }

    body {
      height: 100% !important;
      margin: 0 !important;
      padding: 0 !important;
      width: 100% !important;
    }

    .social-links a {
      text-decoration: none;
    }

    /* MOBILE STYLES */
    @media screen and (max-width: 600px) {
      .intro-table, .logo-table {
        width: 100% !important;
      }

      .logo-table .top-bar {
        border-radius: 0px !important;
      }

      .logo-table .logo {
        margin-left: 5px !important;
        margin-top: 18px !important;
        height: 70px !important;
        width: 75px !important;
      }

      .intro-table td {
        padding: 0px 10px 0px 30px !important;
      }

      .intro-table p {
        margin-bottom: 15px !important;
      }

      .intro-table .shield {
        height: 75px !important;
        width: 81px !important;
        margin-bottom: -30px !important;
      }

      .hidden-row {
        display: none !important;
      }

      .intro-title {
        font-size: 16px !important;
        font-weight: 600 !important;
      }

      .otp-table {
        width: 100% !important;
      }

      .otp-table .otp {
        font-family: 'Lato', Helvetica, Arial, sans-serif;
        font-size: 30px !important;
        font-weight: 900 !important;
        line-height: 23px !important;
      }

      .otp-table tr:first-child td {
        padding: 25px 0px 10px 0px;
      }

      .otp-table .otp-message {
        font-size: 16px !important;
      }

      .otp-table .otp p {
        margin: 0px 0px 8px 0px !important;
      }

      .otp-table .message {
        padding: 20px 30px !important;
        padding-top: 0px;
      }

      .otp-table .message .message-body {
        font-size: 16px;
      }

      .otp-table .regards {
        padding: 0px 30px 50px 30px !important;
        border-radius: 0px 0px 4px 4px;
        color: #666666;
        font-family: 'Lato', Helvetica, Arial, sans-serif;
        font-size: 16px !important;
        font-weight: 400;
        line-height: 25px;
      }

      h1 {
        font-size: 32px !important;
        line-height: 32px !important;
      }

      .socials-table {
        padding-top: 5px !important;
      }

      .socials {
        padding: 5px 30px 0px 30px !important;
      }

      .social-links {
        padding: 10px 30px 5px 30px !important;
        margin: 0px !important;
      }

      .social-links img {
        height: 28px !important;
        width: 30px !important;
      }

      .credits {
        padding: 10px 30px 3px 30px !important;
        font-size: 13px !important;
      }

      .credits p {
        font-size: 13px !important;
      }

      .credits div a {
        margin-right: 5px !important;
      }
    }

    /* ANDROID CENTER FIX */
    div[style*="margin: 16px 0;"] {
      margin: 0 !important;
    }
  </style>
</head>

<body style="
      background-color: #f4f4f4;
      margin: 0 !important;
      padding: 0 !important;
    ">
  <table border="0" cellpadding="0" cellspacing="0" width="100%" height="100%">
    <tr class="hidden-row">
      <td bgcolor="#f4f4f4" valign="center">
        <table border="0" cellpadding="0" cellspacing="0" width="100%" style="max-width: 600px;">
          <tr>
            <td align="center" align="top" style="padding: 40px 10px 40px 10px;"></td>
          </tr>
        </table>
      </td>
    </tr>
    <tr>
      <td bgcolor="#f4f4f4" align="center" style="padding: 0px">
        <table border="0" cellpadding="0" cellspacing="0" width="85%" style="max-width: 600px" class="logo-table">
          <tr>
            <td bgcolor="#66429C" valign="left" class="top-bar" style="
                  padding: 0px 25px 20px 25px;
                  position: relative;
                  border-radius: 4px 4px 0px 0px;
                  color: #fff;
                  font-family: 'Lato', Helvetica, Arial, sans-serif;
                  font-size: 20px;
                  font-weight: 400;
                  line-height: 22px;
                ">
              <img src="https://20261894.fs1.hubspotusercontent-na1.net/hubfs/20261894/assets/white_logo.png"
                width="105" height="100" class="logo" style="margin-left: 18px; margin-top: 20px;" />

            </td>
          </tr>
        </table>
        <table border="0" cellpadding="0" cellspacing="0" width="85%" style="max-width: 600px; border-collapse: collapse !important;" class="intro-table">
          <tr>
            <td bgcolor="#66429C" align="right" class="shield-container" style="
            padding: 10px 25px 10px 25px;
            position: relative;
            font-family: 'Lato', Helvetica, Arial, sans-serif;
            font-size: 20px;
            font-weight: 400;
            border: none;
            line-height: 22px;
          ">
              <img src="https://20261894.fs1.hubspotusercontent-na1.net/hubfs/20261894/assets/shield.png" width="100"
                height="105" class="shield" style="margin-right: 20px;  margin-bottom: -35px"/>
            </td>
          </tr>
        </table>
        <table border="0" cellpadding="0" cellspacing="0" width="85%" style="max-width: 600px; border-collapse: collapse !important;" class="intro-table">
          <tr>
            <td bgcolor="#66429C" align="left" class="top-bar" style="
             padding: 0px 45px 5px 45px;
            color: #fff;
            font-family: 'Lato', Helvetica, Arial, sans-serif;
            font-size: 20px;
            font-weight: 400;
            line-height: 22px;
            border-top-style: none;
          ">
              <p style="" class="intro-title">Hello,</p>
            </td>
          </tr>
        </table>
      </td>
    </tr>
    <tr>
      <td bgcolor="#f4f4f4" align="center" style="padding: 0px">
        <table border="0" cellpadding="0" cellspacing="0" width="85%" style="max-width: 600px" class="otp-table">
          <tr>
            <td bgcolor="#ffffff" align="left" style="
                  padding: 45px 30px 20px 30px;
                  color: #000000;
                  font-family: 'Lato', Helvetica, Arial, sans-serif;
                  font-size: 18px;
                  font-weight: 400;
                  line-height: 25px;
                  text-align: center;
                ">
              <p style="margin: 0" class="otp-message">
                Your Be.Well verification code is:
              </p>
            </td>
          </tr>

          <tr>
            <td bgcolor="#ffffff" align="center" style="
                  color: #000000;
                  font-family: 'Lato', Helvetica, Arial, sans-serif;
                  font-size: 40px;
                  font-weight: 900;
                  line-height: 40px;
                " class="otp">
              <!--Plug in OTP here-->
              <p style="margin-top: 0">{{.}}</p>
            </td>
          </tr>

          <tr>
            <td bgcolor="#ffffff" align="left" style="
                  padding: 10px 45px 50px 45px;
                  border-radius: 0px 0px 4px 4px;
                  color: #666666;
                  font-family: 'Lato', Helvetica, Arial, sans-serif;
                  font-size: 18px;
                  font-weight: 400;
                  line-height: 25px;
                " class="message">
              <p style="margin: 0" class="message-body">
                This code will be active for 60 minutes. If you don't enter it
                on the page you just visited within that time, you can resend
                it from the same page.
              </p>
            </td>
          </tr>

          <tr>
            <td bgcolor="#ffffff" align="left" style="
                  padding: 0px 45px 68px 45px;
                  border-radius: 0px 0px 4px 4px;
                  color: #666666;
                  font-family: 'Lato', Helvetica, Arial, sans-serif;
                  font-size: 18px;
                  font-weight: 400;
                  line-height: 25px;
                " class="regards">
              <p style="margin: 0">
                <span class="regards-message">Kind regards,</span><br />
                <b style="color: #000000;"><span class="team">The Be.Well Team</span></b>
              </p>
            </td>
          </tr>
        </table>
      </td>
    </tr>
    <table border="0" cellpadding="0" cellspacing="0" width="100%" style="padding-top: 40px" class="socials-table">
      <tr>
        <td bgcolor="#f4f4f4" align="center" style="
            padding: 30px 30px 5px 30px;
            border-radius: 0px 0px 4px 4px;
            color: #000000;
            font-family: 'Lato', Helvetica, Arial, sans-serif;
            font-size: 16px;
            font-weight: 400;
            line-height: 25px;
            text-align: center;
          " class="socials">
          <div style="margin: 0" class="social-links">
            <a target="_blank" href="https://www.instagram.com/BeWellBySlade360/" style="margin: 0px  5px;">
              <img
                src="https://20261894.fs1.hubspotusercontent-na1.net/hubfs/20261894/assets/insta-removebg-preview.png"
                width="35" height="30">
            </a>
            <a target="_blank" href="https://www.facebook.com/BeWellbySlade360" style="margin: 0px  5px;">
              <img src="https://20261894.fs1.hubspotusercontent-na1.net/hubfs/20261894/assets/fb-removebg-preview.png"
                width="33" height="27">
            </a>
            <a target="_blank" href="https://www.linkedin.com/showcase/be-well-by-slade-360/" style="margin: 0px  5px;">
              <img src="https://20261894.fs1.hubspotusercontent-na1.net/hubfs/20261894/assets/in-removebg-preview.png"
                width="35" height="30">
            </a>
            <a target="_blank" href="https://www.youtube.com/channel/UCCxjPFz0TxVO6kj-go5vkkA"
              style="margin: 0px  5px;">
              <img src="https://20261894.fs1.hubspotusercontent-na1.net/hubfs/20261894/assets/yt-removebg-preview.png"
                width="35" height="30">
            </a>
            <a target="_blank" href="https://twitter.com/BeWellApp_" style="margin: 0px  5px;">
              <img
                src="https://20261894.fs1.hubspotusercontent-na1.net/hubfs/20261894/assets/tweet-removebg-preview.png"
                width="35" height="30">
            </a>
          </div>
        </td>
      </tr>
      <tr>
        <td bgcolor="#f4f4f4" align="center" style="
              padding: 25px 30px 5px 30px;
              border-radius: 0px 0px 4px 4px;
              color: #000000;
              font-family: 'Lato', Helvetica, Arial, sans-serif;
              font-size: 16px;
              font-weight: 400;
              line-height: 25px;
              text-align: center;
            " class="credits">
          <p style="margin: 0">
            One Padmore Place, George Padmore Road, Nairobi, Kenya
          </p>
        </td>
      </tr>
      <tr>
        <td bgcolor="#f4f4f4" align="center" style="
              padding: 0px 30px 5px 30px;
              border-radius: 0px 0px 4px 4px;
              color: #000000;
              font-family: 'Lato', Helvetica, Arial, sans-serif;
              font-size: 16px;
              font-weight: 400;
              line-height: 25px;
              text-align: center;
            " class="credits">
          <p style="margin: 0">
            If you have any questions or feedback feel free to contact us,
          </p>
        </td>
      </tr>
      <tr>
        <td bgcolor="#f4f4f4" align="center" style="
              padding: 0px 30px 5px 30px;
              border-radius: 0px 0px 4px 4px;
              color: #000000;
              font-family: 'Lato', Helvetica, Arial, sans-serif;
              font-size: 16px;
              font-weight: 400;
              line-height: 25px;
              text-align: center;
            " class="credits">
          <p style="margin: 0">
            <b>feedback@bewell.co.ke</b>
          </p>
        </td>
      </tr>
      <tr>
        <td bgcolor="#f4f4f4" align="center" style="
              padding: 0px 30px 10px 30px;
              border-radius: 0px 0px 4px 4px;
              color: #000000;
              font-family: 'Lato', Helvetica, Arial, sans-serif;
              font-size: 16px;
              font-weight: 400;
              line-height: 25px;
              text-align: center;
            " class="credits">
          <p style="margin: 0">
            <b>0790 360 360</b>
          </p>
        </td>
      </tr>
      <tr>
        <td bgcolor="#f4f4f4" align="center" style="
              padding: 0px 30px 10px 30px;
              border-radius: 0px 0px 4px 4px;
              color: #000000;
              font-family: 'Lato', Helvetica, Arial, sans-serif;
              font-size: 16px;
              font-weight: 400;
              line-height: 25px;
              text-align: center;
            " class="credits">
          <div style="">
            <a target="_blank" href="https://a.bewell.co.ke/loan-terms.html"
              style="margin: 0; text-decoration: none; color: #66429C; margin-right: 20px; font-size: 16px;">
              <b>Terms of use</b>
            </a>
            <a target="_blank" href="https://a.bewell.co.ke/privacy.html"
              style="margin: 0; text-decoration: none; color: #66429C; margin-left: 20px; font-size: 16px;">
              <b>Privacy Policy</b>
            </a>
          </div>
        </td>
      </tr>
      <tr>
        <td bgcolor="#f4f4f4" align="center" style="
                padding: 0px 30px 20px 30px;
                border-radius: 0px 0px 4px 4px;
                color: #000000;
                font-family: 'Lato', Helvetica, Arial, sans-serif;
                font-size: 14px;
                font-weight: 400;
                line-height: 25px;
                text-align: center;
              " class="credits">
          <p>© 2022 All rights reserved Savannah Informatics</p>
          </div>
        </td>
      </tr>
    </table>
    <table border="0" cellpadding="0" cellspacing="0" width="100%" style="padding-top: 20px">
      <tr>
        <td bgcolor="#f4f4f4" align="center" style="
          padding: 20px 30px 20px 30px;
          border-radius: 0px 0px 4px 4px;
          color: #000000;
          font-family: 'Lato', Helvetica, Arial, sans-serif;
          font-size: 18px;
          font-weight: 400;
          line-height: 25px;
        ">
        </td>
      </tr>
    </table>
  </table>
  <script src="https://cdn.jsdelivr.net/npm/publicalbum@latest/embed-ui.min.js" async></script>
  <!-- The core Firebase JS SDK is always required and must be listed first -->
  <script src="https://www.gstatic.com/firebasejs/8.7.1/firebase-app.js"></script>

  <script src="https://www.gstatic.com/firebasejs/8.7.1/firebase-analytics.js"></script>

  <script>
    var firebaseConfig = {
      apiKey: "AIzaSyAv2aRsSSHkOR6xGwwaw6-UTkvED3RNlBQ",
      authDomain: "bewell-app.firebaseapp.com",
      databaseURL: "https://bewell-app.firebaseio.com",
      projectId: "bewell-app",
      storageBucket: "bewell-app.appspot.com",
      messagingSenderId: "841947754847",
      appId: "1:841947754847:web:6304157d32c82fd96686ea",
      measurementId: "G-6XTZEB5070",
    };
    firebase.initializeApp(firebaseConfig);
    const analytics = firebase.analytics();

    analytics.logEvent("opened_otp_email");
  </script>
</body>

</html>

`
